package zigbee

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
)

type Z2MManager struct {
	config  *config.Zigbee2Mqtt
	handler map[string]Z2MMessageHandler
	client  mqtt.Client
}

type Z2MMessageHandler func(mqtt.Message)

func NewZ2MManager(config *config.Zigbee2Mqtt) *Z2MManager {
	return &Z2MManager{
		config:  config,
		handler: make(map[string]Z2MMessageHandler),
	}
}

func (z2m *Z2MManager) Connect() {
	misc.Log.Infof("Connecting to zigbee broker: %v", z2m.config.Connection)
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
		panic(token.Error())
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

	misc.Log.Debugf("Registered handler for %v", topicName)
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
		misc.Log.Warnf("Could not marshall message for %v: %v", topicName, payload)
		return
	}

	z2m.client.Publish(topicName, 1, false, string(data))
	misc.Log.Debugf("Send to %v: %v", topicName, string(data))
}

func (z2m *Z2MManager) messageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		handler, ok := z2m.handler[msg.Topic()]
		if ok && handler != nil {
			handler(msg)
		} else {
			misc.Log.Warnf("Could ot find handler for message in topic %s", msg.Topic())
		}
	}
}

func (z2m *Z2MManager) onConnectHandler() mqtt.OnConnectHandler {
	return func(client mqtt.Client) {
		misc.Log.Infof("Connected to zigbee broker")
	}
}

func (z2m *Z2MManager) onConnectionLostHandler() mqtt.ConnectionLostHandler {
	return func(client mqtt.Client, err error) {
		misc.Log.Warnf("Connection to zigbee broker lost!: %v", err)
		z2m.Connect()

		if len(z2m.handler) > 0 {
			misc.Log.Debugf("There were message handlers registered before connection to broker has been established.")
			for t, h := range z2m.handler {
				z2m.Subscribe(t, h)
			}
		}
	}
}
