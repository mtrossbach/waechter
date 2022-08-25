package keypad

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/system"
)

type keypad struct {
	deviceInfo    model2.Z2MDeviceInfo
	connector     *connector.Connector
	systemControl system.Controller
	targetTopic   string
}

func New(deviceInfo model2.Z2MDeviceInfo, connector *connector.Connector) *keypad {
	return &keypad{
		deviceInfo:  deviceInfo,
		connector:   connector,
		targetTopic: fmt.Sprintf("%v/set", deviceInfo.FriendlyName),
	}
}

func (s *keypad) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *keypad) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *keypad) GetSubsystem() string {
	return model2.SubsystemName
}

func (s *keypad) GetType() system.DeviceType {
	return system.Keypad
}

func (s *keypad) OnSystemStateChanged(state system.State, aMode system.ArmingMode, aType system.AlarmType) {
	s.sendState()
}

func (s *keypad) OnDeviceAnnounced() {
	s.sendState()
}

func (s *keypad) Setup(systemControl system.Controller) {
	s.systemControl = systemControl
	s.connector.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
	s.sendState()
}

func (s *keypad) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *keypad) handleMessage(msg mqtt.Message) {
	var payload statusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if payload.Battery > 0 {
		level := float32(payload.Battery) / float32(100)
		s.systemControl.ReportBatteryLevel(level, s)
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.Linkquality)/float32(255), s)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s)
	}

	if len(payload.Action) > 0 {
		s._sendState(payload.Action, &payload.ActionTransaction) //Send confirmation (required for some devices)

		if payload.Action == "arm_day_zones" {
			s.systemControl.Arm(system.StayMode, s)
		} else if payload.Action == "arm_all_zones" {
			s.systemControl.Arm(system.AwayMode, s)
		} else if payload.Action == "disarm" {
			s.systemControl.Disarm(payload.ActionCode, s)
		} else if payload.Action == "panic" {
			s.systemControl.Alarm(system.PanicAlarm, s)
		}
	}
}

func (s *keypad) _sendState(state string, transactionId *int) {
	payload := newStatePayload(state, transactionId)
	s.connector.Publish(s.targetTopic, payload)
}

func (s *keypad) sendState() {
	if s.systemControl == nil {
		return
	}
	s._sendState(systemStateToDeviceState(s.systemControl.GetState(), s.systemControl.GetArmingMode(), s.systemControl.GetAlarmType()), nil)

}
