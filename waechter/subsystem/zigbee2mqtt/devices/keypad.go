package devices

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/waechter/misc"
	"github.com/mtrossbach/waechter/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/waechter/system"
)

type setStatePayload struct {
	ArmMode ArmMode `json:"arm_mode"`
}
type ArmMode struct {
	Mode string `json:"mode"`
}
type genericKeypad struct {
	zdevice       model.ZigbeeDevice
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
	targetTopic   string
}

func newGenericKeypad(zdevice model.ZigbeeDevice, z2mManager *zigbee.Z2MManager) *genericKeypad {
	return &genericKeypad{
		zdevice:     zdevice,
		z2mManager:  z2mManager,
		targetTopic: fmt.Sprintf("%v/set", zdevice.FriendlyName),
	}
}

func (s *genericKeypad) GetId() string {
	return s.zdevice.IeeeAddress
}

func (s *genericKeypad) GetDisplayName() string {
	return s.zdevice.FriendlyName
}

func (s *genericKeypad) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericKeypad) GetType() system.DeviceType {
	return system.Keypad
}

func (s *genericKeypad) OnSystemStateChanged(state system.State) {
	misc.Log.Debugf("State changed to %v", state)
	s.sendState()
}

func (s *genericKeypad) Setup(systemControl system.SystemControl) {
	misc.Log.Debugf("Setup device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.zdevice.FriendlyName, s.handleMessage)
	s.sendState()
}

func (s *genericKeypad) Teardown() {
	misc.Log.Debugf("Teardown device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.zdevice.FriendlyName)
}

func (s *genericKeypad) handleMessage(msg mqtt.Message) {
	misc.Log.Debugf("Got data: %v", string(msg.Payload()))
}

func (s *genericKeypad) systemStateToDeviceState(state system.State) string {
	if state == system.Disarmed {
		return "disarm"
	} else if state == system.ArmingStay {
		return "arming_stay"
	} else if state == system.ArmingAway {
		return "arming_away"
	} else if state == system.ArmedStay {
		return "arm_night_zones"
	} else if state == system.ArmedAway {
		return "arm_all_zones"
	} else if state == system.InAlarm {
		return "in_alarm"
	} else if state == system.EntryDelay {
		return "entry_delay"
	}

	return ""
}

func (s *genericKeypad) sendState() {
	if s.systemControl == nil {
		return
	}
	payload := setStatePayload{
		ArmMode: ArmMode{
			Mode: s.systemStateToDeviceState(s.systemControl.GetState()),
		},
	}

	s.z2mManager.Publish(s.targetTopic, payload)
}
