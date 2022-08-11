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

type Z2MManager struct {
	config  *config.Zigbee2Mqtt
	handler map[string]Z2MMessageHandler
	client  mqtt.Client
	log     zerolog.Logger
}

type Z2MMessageHandler func(mqtt.Message)

func NewZ2MManager(config *config.Zigbee2Mqtt) *Z2MManager {
	return &Z2MManager{
		config:  config,
		handler: make(map[string]Z2MMessageHandler),
		log:     misc.Logger("Z2MManager"),
	}
}

func (z2m *Z2MManager) Connect() {
	z2m.log.Info().Str("connection", z2m.config.Connection).Str("clientId", z2m.config.ClientId).Str("username", z2m.config.Username).Msg("Connecting to mqtt broker...")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(z2m.config.Connection)
	opts.SetClientID(z2m.config.ClientId)
	opts.SetUsername(z2m.config.Username)
	opts.SetPassword(z2m.config.Password)
	opts.SetDefaultPublishHandler(z2m.messageHandler())

	opts.OnConnect = z2m.onConnectHandler()
	opts.OnConnectionLost = z2m.onConnectionLostHandler()
	z2m.client = mqtt.NewClient(opts)
	if token := z2m.client.Connect(); token.Wait() && token.Error() != nil {
		z2m.log.Error().Err(token.Error()).Msg("Could not connect to mqtt broker")
	}
}

func (z2m *Z2MManager) Disconnect() {
	z2m.client.Disconnect(100)
}

func (z2m *Z2MManager) Subscribe(topic string, handler Z2MMessageHandler) {
	topicName := fmt.Sprintf("%s/%s", z2m.config.BaseTopic, topic)
	if strings.HasPrefix(topic, z2m.config.BaseTopic) {
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

func (z2m *Z2MManager) Unsubscribe(topic string) {
	topicName := fmt.Sprintf("%s/%s", z2m.config.BaseTopic, topic)
	z2m.client.Unsubscribe(topicName)
	delete(z2m.handler, topicName)
}

func (z2m *Z2MManager) Publish(topic string, payload interface{}) {
	topicName := fmt.Sprintf("%s/%s", z2m.config.BaseTopic, topic)
	data, err := json.Marshal(payload)
	if err != nil {
		z2m.log.Error().Str("topic", topicName).Interface("payload", payload).Msg("Could not parse payload")
		return
	}

	z2m.client.Publish(topicName, 1, false, string(data))
	z2m.log.Debug().Str("topic", topicName).RawJSON("msg", data).Msg("Sent message.")
}

func (z2m *Z2MManager) messageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		handler, ok := z2m.handler[msg.Topic()]
		if ok && handler != nil {
			handler(msg)
		} else {
			z2m.log.Warn().Str("topic", msg.Topic()).Msg("Could not find handler for message.")
		}
	}
}

func (z2m *Z2MManager) onConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		z2m.log.Info().Msg("Connected to mqtt broker")
	}
}

func (z2m *Z2MManager) onConnectionLostHandler() mqtt.ConnectionLostHandler {
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
