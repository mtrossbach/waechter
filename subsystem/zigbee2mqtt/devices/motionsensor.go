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

type motionsensorStatusPayload struct {
	Battery     int  `json:"battery"`
	Linkquality int  `json:"linkquality"`
	Occupancy   bool `json:"occupancy"`
	Tamper      bool `json:"tamper"`
}

type genericMotionSensor struct {
	deviceInfo    model.Z2MDeviceInfo
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
	log           zerolog.Logger
}

func newGenericMotionSensor(deviceInfo model.Z2MDeviceInfo, z2mManager *zigbee.Z2MManager) *genericMotionSensor {
	return &genericMotionSensor{
		deviceInfo: deviceInfo,
		z2mManager: z2mManager,
		log:        misc.Logger("Z2MMotionSensor"),
	}
}

func (s *genericMotionSensor) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *genericMotionSensor) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *genericMotionSensor) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericMotionSensor) GetType() system.DeviceType {
	return system.MotionSensor
}

func (s *genericMotionSensor) OnSystemStateChanged(state system.State) {

}

func (s *genericMotionSensor) OnDeviceAnnounced() {

}

func (s *genericMotionSensor) Setup(systemControl system.SystemControl) {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Setup device")
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *genericMotionSensor) Teardown() {
	s.log.Debug().Str("type", string(s.GetType())).Str("id", s.GetId()).Str("displayName", s.GetDisplayName()).Msg("Tear down device")
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *genericMotionSensor) handleMessage(msg mqtt.Message) {
	var payload motionsensorStatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.log.Warn().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	s.log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBattery(s, float32(payload.Battery)/float32(100))
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportBattery(s, float32(payload.Linkquality)/float32(255))
	}

	if payload.Tamper {
		s.systemControl.ReportTampered(s)
	}

	if payload.Occupancy {
		s.systemControl.ReportTriggered(s)
	}
}
