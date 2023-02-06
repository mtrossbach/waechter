package zigbee2mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type connection struct {
	config  config.Zigbee2MqttConfig
	handler map[string]MessageHandler
	client  mqtt.Client

	OnConnect        ConnectedHandler
	OnConnectionLost ConnectionLostHandler
}

type MessageHandler func(mqtt.Message)

type ConnectedHandler func(conn *connection)
type ConnectionLostHandler func(conn *connection, err error)

func newConnection(configuration config.Zigbee2MqttConfig) *connection {
	return &connection{
		config:  configuration,
		handler: make(map[string]MessageHandler),
	}
}

func (c *connection) Connect() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.Url)
	opts.SetClientID(c.config.ClientId)
	opts.SetUsername(c.config.Username)
	opts.SetPassword(c.config.Password)
	opts.SetDefaultPublishHandler(c.messageHandler())

	opts.OnConnect = c.onConnectHandler()
	opts.OnConnectionLost = c.onConnectionLostHandler()
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		if c.OnConnectionLost != nil {
			c.OnConnectionLost(c, token.Error())
		}
		go c.reconnect()
	}
}

func (c *connection) Disconnect() {
	for k := range c.handler {
		c.Unsubscribe(k)
	}
	c.handler = map[string]MessageHandler{}

	if c.OnConnectionLost != nil {
		c.OnConnectionLost(c, nil)
	}
	c.client.Disconnect(100)
}

func (c *connection) Subscribe(topic string, handler MessageHandler) bool {
	topicName := fmt.Sprintf("%s/%s", c.config.BaseTopic, topic)
	if strings.HasPrefix(topic, c.config.BaseTopic) {
		topicName = topic
	}

	c.handler[topicName] = handler
	if c.client == nil {
		return false
	}
	token := c.client.Subscribe(topicName, 1, nil)
	ok := token.WaitTimeout(10 * time.Second)
	if token.Error() != nil || !ok {
		log.Error().Str("topic", topicName).Err(token.Error()).Msg("Could not register handler")
		return false
	} else {
		log.Debug().Str("topic", topicName).Msg("Registered handler")
		return true
	}
}

func (c *connection) Unsubscribe(topic string) {
	topicName := fmt.Sprintf("%s/%s", c.config.BaseTopic, topic)
	c.client.Unsubscribe(topicName)
	delete(c.handler, topicName)
}

func (c *connection) Publish(topic string, payload interface{}) {
	topicName := fmt.Sprintf("%s/%s", c.config.BaseTopic, topic)
	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Str("topic", topicName).Interface("payload", payload).Msg("Could not parse payload")
		return
	}

	c.client.Publish(topicName, 1, false, string(data))
	log.Debug().Str("topic", topicName).RawJSON("msg", data).Msg("Sent message.")
}

func (c *connection) messageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		handler, ok := c.handler[msg.Topic()]
		if ok && handler != nil {
			go handler(msg)
		} else {
			log.Error().Str("topic", msg.Topic()).Msg("Could not find handler for message.")
		}
	}
}

func (c *connection) onConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		for topic, handler := range c.handler {
			c.Subscribe(topic, handler)
		}
		if c.OnConnect != nil {
			c.OnConnect(c)
		}
	}
}

func (c *connection) onConnectionLostHandler() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		if c.OnConnectionLost != nil {
			c.OnConnectionLost(c, err)
		}
		go c.reconnect()
	}
}

func (c *connection) reconnect() {
	<-time.After(1 * time.Second)
	c.Connect()
}
