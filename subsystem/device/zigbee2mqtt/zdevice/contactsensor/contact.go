package contactsensor

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/system"
)

type contactSensor struct {
	system.Device
	connector     *connector.Connector
	systemControl system.Controller
	readTopic     string
}

func New(device system.Device) *contactSensor {
	return &contactSensor{
		Device: system.Device{
			Id:   device.Id,
			Name: device.Name,
			Type: system.MotionSensor,
		},
		readTopic: device.Name,
	}
}

func (s *contactSensor) OnDeviceAnnounced() {

}

func (s *contactSensor) UpdateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType) {

}

func (s *contactSensor) Setup(connector *connector.Connector, systemControl system.Controller) {
	s.systemControl = systemControl
	s.connector = connector
	s.connector.Subscribe(s.readTopic, s.handleMessage)
	system.DInfo(s.Device).Msg("Activated.")
}

func (s *contactSensor) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.readTopic)
	s.connector = nil
	system.DInfo(s.Device).Msg("Deactivated.")
}

func (s *contactSensor) handleMessage(msg mqtt.Message) {
	var payload statusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBatteryLevel(float32(payload.Battery)/float32(100), s.Device)
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.Linkquality)/float32(255), s.Device)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s.Device)
	}

	if !payload.Contact {
		s.systemControl.Alarm(system.BurglarAlarm, s.Device)
	}
}
