package zdevice

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
)

type siren struct {
	system.Device
	connector     *connector.Connector
	systemControl device.SystemController
	writeTopic    string
	readTopic     string
}

func NewSiren(device system.Device) *siren {
	return &siren{
		Device: system.Device{
			Id:   device.Id,
			Name: device.Name,
			Type: system.Siren,
		},
		readTopic:  device.Name,
		writeTopic: fmt.Sprintf("%v/set", device.Name),
	}
}

func (s *siren) UpdateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType) {
	s.sendState()
}

func (s *siren) Setup(connector *connector.Connector, systemControl device.SystemController) {
	s.systemControl = systemControl
	s.connector = connector
	s.connector.Subscribe(s.readTopic, s.handleMessage)
	system.DInfo(s.Device).Msg("Activated.")
}

func (s *siren) OnDeviceAnnounced() {
	s.sendState()
}

func (s *siren) sendState() {
	var payload sirenWarning
	if s.systemControl.GetAlarmType() != system.NoAlarm && cfg.GetBool(cEnabled) {
		payload = newSirenWarningPayload(s.systemControl.GetAlarmType())
	} else {
		payload = newSirenWarningPayload(system.NoAlarm)
	}
	s.connector.Publish(s.writeTopic, payload)
}

func (s *siren) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.readTopic)
	s.connector = nil
	system.DInfo(s.Device).Msg("Deactivated.")
}

func (s *siren) handleMessage(msg mqtt.Message) {
	var payload sirenStatus
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBatteryLevel(float32(payload.Battery)/float32(100), s.Device)
	}

	if payload.LinkQuality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.LinkQuality)/float32(255), s.Device)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s.Device)
	}
}
