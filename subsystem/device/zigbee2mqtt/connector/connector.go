package connector

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
	"github.com/rs/zerolog"
)

type Connector struct {
	handler map[string]Z2MMessageHandler
	client  mqtt.Client
	log     zerolog.Logger
}

type Z2MMessageHandler func(mqtt.Message)

func New() *Connector {
	return &Connector{
		handler: make(map[string]Z2MMessageHandler),
		log:     misc.Logger("Z2MConnector"),
	}
}

func (z2m *Connector) Connect() {
	z2m.log.Info().Str("connection", config.GetString(cConnection)).Str("clientId", config.GetString(cClientId)).Str("username", config.GetString(cUsername)).Msg("Connecting to mqtt broker...")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.GetString(cConnection))
	opts.SetClientID(config.GetString(cClientId))
	opts.SetUsername(config.GetString(cUsername))
	opts.SetPassword(config.GetString(cPassword))
	opts.SetDefaultPublishHandler(z2m.messageHandler())

	opts.OnConnect = z2m.onConnectHandler()
	opts.OnConnectionLost = z2m.onConnectionLostHandler()
	z2m.client = mqtt.NewClient(opts)
	if token := z2m.client.Connect(); token.Wait() && token.Error() != nil {
		z2m.log.Error().Err(token.Error()).Msg("Could not connect to mqtt broker")
	}
}

func (z2m *Connector) Disconnect() {
	z2m.client.Disconnect(100)
}

func (z2m *Connector) Subscribe(topic string, handler Z2MMessageHandler) {
	topicName := fmt.Sprintf("%s/%s", config.GetString(cBaseTopic), topic)
	if strings.HasPrefix(topic, config.GetString(cBaseTopic)) {
		topicName = topic
	}

	z2m.handler[topicName] = handler
	token := z2m.client.Subscribe(topicName, 1, nil)
	token.Wait()
	if token.Error() != nil {
		z2m.log.Error().Str("topic", topicName).Err(token.Error()).Msg("Could not register handler")
	} else {
		z2m.log.Debug().Str("topic", topicName).Msg("Registered handler")
	}
}

func (z2m *Connector) Unsubscribe(topic string) {
	topicName := fmt.Sprintf("%s/%s", config.GetString(cBaseTopic), topic)
	z2m.client.Unsubscribe(topicName)
	delete(z2m.handler, topicName)
}

func (z2m *Connector) Publish(topic string, payload interface{}) {
	topicName := fmt.Sprintf("%s/%s", config.GetString(cBaseTopic), topic)
	data, err := json.Marshal(payload)
	if err != nil {
		z2m.log.Error().Str("topic", topicName).Interface("payload", payload).Msg("Could not parse payload")
		return
	}

	z2m.client.Publish(topicName, 1, false, string(data))
	z2m.log.Debug().Str("topic", topicName).RawJSON("msg", data).Msg("Sent message.")
}

func (z2m *Connector) messageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		handler, ok := z2m.handler[msg.Topic()]
		if ok && handler != nil {
			handler(msg)
		} else {
			z2m.log.Warn().Str("topic", msg.Topic()).Msg("Could not find handler for message.")
		}
	}
}

func (z2m *Connector) onConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		z2m.log.Info().Msg("Connected to mqtt broker")
	}
}

func (z2m *Connector) onConnectionLostHandler() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		z2m.log.Error().Err(err).Msg("Connection to mqtt broker lost!")
		z2m.Connect()

		if len(z2m.handler) > 0 {
			z2m.log.Debug().Msg("There were message handlers registered before connection to mqtt broker has been established.")
			for t, h := range z2m.handler {
				z2m.Subscribe(t, h)
			}
		}
	}
}