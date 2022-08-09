package devices

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type contactsensorStatusPayload struct {
	Battery     int  `json:"battery"`
	Linkquality int  `json:"linkquality"`
	Contact     bool `json:"contact"`
	Tamper      bool `json:"tamper"`
}

type genericContactSensor struct {
	deviceInfo    model.Z2MDeviceInfo
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
	log           zerolog.Logger
}

func newGenericContactSensor(deviceInfo model.Z2MDeviceInfo, z2mManager *zigbee.Z2MManager) *genericContactSensor {
	return &genericContactSensor{
		deviceInfo: deviceInfo,
		z2mManager: z2mManager,
		log:        misc.Logger("Z2MContactSensor"),
	}
}

func (s *genericContactSensor) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *genericContactSensor) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *genericContactSensor) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericContactSensor) GetType() system.DeviceType {
	return system.ContactSensor
}

func (s *genericContactSensor) OnDeviceAnnounced() {

}

func (s *genericContactSensor) OnSystemStateChanged(state system.State) {

}

func (s *genericContactSensor) Setup(systemControl system.SystemControl) {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Setup device")
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *genericContactSensor) Teardown() {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Tear down device")
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *genericContactSensor) handleMessage(msg mqtt.Message) {
	var payload contactsensorStatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.log.Warn().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	s.log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

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
