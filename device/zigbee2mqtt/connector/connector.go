package connector

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Connector struct {
	handler             map[string]MessageHandler
	client              mqtt.Client
	baseTopic           string
	disconnectedHandler DisconnectedHandler
	connectedHandler    ConnectedHandler
}

type MessageHandler func(mqtt.Message)

type ConnectedHandler func()
type DisconnectedHandler func(err error)

func New() *Connector {
	return &Connector{
		handler: make(map[string]MessageHandler),
	}
}

type Options struct {
	Uri       string
	ClientId  string
	Username  string
	Password  string
	BaseTopic string
}

func (c *Connector) Connect(options Options, connectedHandler ConnectedHandler, disconnectedHandler DisconnectedHandler) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(options.Uri)
	opts.SetClientID(options.ClientId)
	opts.SetUsername(options.Username)
	opts.SetPassword(options.Password)
	c.baseTopic = options.BaseTopic
	c.connectedHandler = connectedHandler
	c.disconnectedHandler = disconnectedHandler
	opts.SetDefaultPublishHandler(c.messageHandler())

	opts.OnConnect = c.onConnectHandler()
	opts.OnConnectionLost = c.onConnectionLostHandler()
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {

		return token.Error()
	}
	return nil
}

func (c *Connector) Disconnect() {
	for k := range c.handler {
		c.Unsubscribe(k)
	}
	c.handler = map[string]MessageHandler{}
	c.connectedHandler = nil
	if c.disconnectedHandler != nil {
		c.disconnectedHandler(nil)
		c.disconnectedHandler = nil
	}
	c.baseTopic = ""
	c.client.Disconnect(100)
}

func (c *Connector) Subscribe(topic string, handler MessageHandler) {
	topicName := fmt.Sprintf("%s/%s", c.baseTopic, topic)
	if strings.HasPrefix(topic, c.baseTopic) {
		topicName = topic
	}

	c.handler[topicName] = handler
	token := c.client.Subscribe(topicName, 1, nil)
	token.Wait()
	if token.Error() != nil {
		log.Error().Str("topic", topicName).Err(token.Error()).Msg("Could not register handler")
	} else {
		log.Debug().Str("topic", topicName).Msg("Registered handler")
	}
}

func (c *Connector) Unsubscribe(topic string) {
	topicName := fmt.Sprintf("%s/%s", c.baseTopic, topic)
	c.client.Unsubscribe(topicName)
	delete(c.handler, topicName)
}

func (c *Connector) Publish(topic string, payload interface{}) {
	topicName := fmt.Sprintf("%s/%s", c.baseTopic, topic)
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Str("topic", topicName).Interface("payload", payload).Msg("Could not parse payload")
		return
	}

	c.client.Publish(topicName, 1, false, string(data))
	log.Debug().Str("topic", topicName).RawJSON("msg", data).Msg("Sent message.")
}

func (c *Connector) messageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		handler, ok := c.handler[msg.Topic()]
		if ok && handler != nil {
			handler(msg)
		} else {
			log.Error().Str("topic", msg.Topic()).Msg("Could not find handler for message.")
		}
	}
}

func (c *Connector) onConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		if c.connectedHandler != nil {
			c.connectedHandler()
		}
	}
}

func (c *Connector) onConnectionLostHandler() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		if c.disconnectedHandler != nil {
			c.disconnectedHandler(err)
		}
	}
}
