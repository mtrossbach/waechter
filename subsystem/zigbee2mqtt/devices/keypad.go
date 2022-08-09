package devices

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type setStatePayload struct {
	ArmMode ArmMode `json:"arm_mode"`
}
type ArmMode struct {
	Mode        string `json:"mode"`
	Transaction *int   `json:"transaction,omitempty"`
}

type keypadStatusPayload struct {
	Action            string `json:"action"`
	ActionCode        string `json:"action_code"`
	ActionTransaction int    `json:"action_transaction"`
	ActionZone        int    `json:"action_zone"`
	Battery           int    `json:"battery"`
	Linkquality       int    `json:"linkquality"`
	Occupancy         bool   `json:"occupancy"`
	Tamper            bool   `json:"tamper"`
}

type genericKeypad struct {
	deviceInfo    model.Z2MDeviceInfo
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
	targetTopic   string
	log           zerolog.Logger
}

func newGenericKeypad(deviceInfo model.Z2MDeviceInfo, z2mManager *zigbee.Z2MManager) *genericKeypad {
	return &genericKeypad{
		deviceInfo:  deviceInfo,
		z2mManager:  z2mManager,
		targetTopic: fmt.Sprintf("%v/set", deviceInfo.FriendlyName),
		log:         misc.Logger("Z2MKeypad"),
	}
}

func (s *genericKeypad) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *genericKeypad) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *genericKeypad) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericKeypad) GetType() system.DeviceType {
	return system.Keypad
}

func (s *genericKeypad) OnSystemStateChanged(state system.State) {
	s.sendState()
}

func (s *genericKeypad) OnDeviceAnnounced() {
	s.sendState()
}

func (s *genericKeypad) Setup(systemControl system.SystemControl) {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Setup device")
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
	s.sendState()
}

func (s *genericKeypad) Teardown() {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Tear down device")
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *genericKeypad) handleMessage(msg mqtt.Message) {
	var payload keypadStatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.log.Warn().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	s.log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if len(payload.Action) > 0 {
		s._sendState(payload.Action, &payload.ActionTransaction)
		if payload.Action == "arm_day_zones" {
			s.systemControl.ArmStay()
		} else if payload.Action == "arm_all_zones" {
			s.systemControl.ArmAway()
		} else if payload.Action == "disarm" {
			s.systemControl.Disarm(payload.ActionCode)
		} else if payload.Action == "panic" {
			s.systemControl.Panic()
		}
	}

	/*
		{"ac_status":false,"action":"arm_all_zones","action_code":"1337",
		"action_transaction":15,"action_zone":23,"battery":100,
		"battery_low":false,"linkquality":54,"occupancy":true,"restore_reports":true,
		"smoke":false,"supervision_reports":true,"tamper":true,"test":false,"trouble":false,"voltage":3100}
	*/
}

func (s *genericKeypad) systemStateToDeviceState(state system.State) string {
	if state == system.Disarmed {
		return "disarm"
	} else if state == system.ArmingStay {
		return "exit_delay"
	} else if state == system.ArmingAway {
		return "exit_delay"
	} else if state == system.ArmedStay {
		return "arm_day_zones"
	} else if state == system.ArmedAway {
		return "arm_all_zones"
	} else if state == system.InAlarm {
		return "in_alarm"
	} else if state == system.EntryDelay {
		return "entry_delay"
	} else if state == system.Panic {
		return "panic"
	}

	return ""
}

func (s *genericKeypad) _sendState(state string, transactionId *int) {
	payload := setStatePayload{
		ArmMode: ArmMode{
			Mode:        state,
			Transaction: transactionId,
		},
	}

	s.z2mManager.Publish(s.targetTopic, payload)
}

func (s *genericKeypad) sendState() {
	if s.systemControl == nil {
		return
	}
	s._sendState(s.systemStateToDeviceState(s.systemControl.GetState()), nil)

}
