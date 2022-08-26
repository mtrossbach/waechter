package zdevice

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
)

type keypad struct {
	system.Device
	connector     *connector.Connector
	systemControl device.SystemController
	writeTopic    string
	readTopic     string
}

func NewKeypad(device system.Device) *keypad {
	return &keypad{
		Device: system.Device{
			Id:   device.Id,
			Name: device.Name,
			Type: system.Keypad,
		},
		readTopic:  device.Name,
		writeTopic: fmt.Sprintf("%v/set", device.Name),
	}
}

func (s *keypad) UpdateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType) {
	s.sendState()
}

func (s *keypad) Setup(connector *connector.Connector, systemControl device.SystemController) {
	s.systemControl = systemControl
	s.connector = connector
	s.connector.Subscribe(s.readTopic, s.handleMessage)
	system.DInfo(s.Device).Msg("Activated.")
	s.sendState()
}
func (s *keypad) OnDeviceAnnounced() {
	s.sendState()
}

func (s *keypad) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.readTopic)
	system.DInfo(s.Device).Msg("Deactivated.")
}

func (s *keypad) handleMessage(msg mqtt.Message) {
	var payload keypadStatus
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if payload.Battery > 0 {
		level := float32(payload.Battery) / float32(100)
		s.systemControl.ReportBatteryLevel(level, s.Device)
	}

	if payload.LinkQuality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.LinkQuality)/float32(255), s.Device)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s.Device)
	}

	if len(payload.Action) > 0 {
		s._sendState(payload.Action, &payload.ActionTransaction) //Send confirmation (required for some devices)

		if payload.Action == "arm_day_zones" {
			s.systemControl.Arm(system.StayMode, s.Device)
		} else if payload.Action == "arm_all_zones" {
			s.systemControl.Arm(system.AwayMode, s.Device)
		} else if payload.Action == "disarm" {
			s.systemControl.Disarm(payload.ActionCode, s.Device)
		} else if payload.Action == "panic" {
			s.systemControl.Alarm(system.PanicAlarm, s.Device)
		}
	}
}

func (s *keypad) _sendState(state string, transactionId *int) {
	payload := newKeypadSetState(state, transactionId)
	s.connector.Publish(s.writeTopic, payload)
}

func (s *keypad) sendState() {
	if s.systemControl == nil {
		return
	}
	s._sendState(sysStateToKeypadState(s.systemControl.GetState(), s.systemControl.GetArmingMode(), s.systemControl.GetAlarmType()), nil)

}
