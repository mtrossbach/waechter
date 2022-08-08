package devices

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/waechter/misc"
	"github.com/mtrossbach/waechter/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/waechter/system"
)

type contactsensorStatusPayload struct {
	Battery     int  `json:"battery"`
	Linkquality int  `json:"linkquality"`
	Contact     bool `json:"contact"`
	Tamper      bool `json:"tamper"`
}

type genericContactSensor struct {
	zdevice       model.ZigbeeDevice
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
}

func newGenericContactSensor(zdevice model.ZigbeeDevice, z2mManager *zigbee.Z2MManager) *genericContactSensor {
	return &genericContactSensor{
		zdevice:    zdevice,
		z2mManager: z2mManager,
	}
}

func (s *genericContactSensor) GetId() string {
	return s.zdevice.IeeeAddress
}

func (s *genericContactSensor) GetDisplayName() string {
	return s.zdevice.FriendlyName
}

func (s *genericContactSensor) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericContactSensor) GetType() system.DeviceType {
	return system.ContactSensor
}

func (s *genericContactSensor) OnSystemStateChanged(state system.State) {
	misc.Log.Debugf("State changed to %v", state)
}

func (s *genericContactSensor) Setup(systemControl system.SystemControl) {
	misc.Log.Debugf("Setup device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.zdevice.FriendlyName, s.handleMessage)
}

func (s *genericContactSensor) Teardown() {
	misc.Log.Debugf("Teardown device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.zdevice.FriendlyName)
}

func (s *genericContactSensor) handleMessage(msg mqtt.Message) {
	var payload contactsensorStatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		misc.Log.Warnf("Could not parse payload: %v", string(msg.Payload()))
		return
	}

	misc.Log.Debugf("Got data: %v", string(msg.Payload()))

	if payload.Battery > 0 {
		s.systemControl.ReportBattery(s, float32(payload.Battery)/float32(100))
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportBattery(s, float32(payload.Linkquality)/float32(255))
	}

	if payload.Tamper {
		s.systemControl.ReportTampered(s)
	}

	if !payload.Contact {
		s.systemControl.ReportTriggered(s)
	}
}
