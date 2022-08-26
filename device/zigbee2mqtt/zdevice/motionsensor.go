package zdevice

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
)

type motionSensor struct {
	system.Device
	connector     *connector.Connector
	systemControl device.SystemController
	readTopic     string
}

func NewMotionSensor(device system.Device) *motionSensor {
	return &motionSensor{
		Device: system.Device{
			Id:   device.Id,
			Name: device.Name,
			Type: system.MotionSensor,
		},
		readTopic: device.Name,
	}
}

func (s *motionSensor) OnDeviceAnnounced() {

}

func (s *motionSensor) UpdateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType) {
}

func (s *motionSensor) Setup(connector *connector.Connector, systemControl device.SystemController) {
	s.systemControl = systemControl
	s.connector = connector
	s.connector.Subscribe(s.readTopic, s.handleMessage)
	system.DInfo(s.Device).Msg("Activated.")
}

func (s *motionSensor) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.readTopic)
	s.connector = nil
	system.DInfo(s.Device).Msg("Deactivated.")
}

func (s *motionSensor) handleMessage(msg mqtt.Message) {
	var payload motionStatus
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBatteryLevel(float32(payload.Battery)/float32(100), s.Device)
	}

	if payload.LinkQuality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.LinkQuality)/float32(255), s.Device)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s.Device)
	}

	if payload.Occupancy {
		s.systemControl.Alarm(system.BurglarAlarm, s.Device)
	}
}
