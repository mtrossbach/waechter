package motionsensor

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type motionSensor struct {
	deviceInfo    model2.Z2MDeviceInfo
	z2mManager    *connector.Z2MManager
	systemControl system.Controller
	log           zerolog.Logger
}

func NewMotionSensor(deviceInfo model2.Z2MDeviceInfo, z2mManager *connector.Z2MManager) *motionSensor {
	return &motionSensor{
		deviceInfo: deviceInfo,
		z2mManager: z2mManager,
		log:        misc.Logger("Z2MMotionSensor"),
	}
}

func (s *motionSensor) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *motionSensor) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *motionSensor) GetSubsystem() string {
	return model2.SubsystemName
}

func (s *motionSensor) GetType() system.DeviceType {
	return system.MotionSensor
}

func (s *motionSensor) OnSystemStateChanged(state system.State, aMode system.ArmingMode, aType system.AlarmType) {

}

func (s *motionSensor) OnDeviceAnnounced() {

}

func (s *motionSensor) Setup(systemControl system.Controller) {
	system.DevLog(s, s.log.Debug()).Msg("Setup zdevice")
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *motionSensor) Teardown() {
	system.DevLog(s, s.log.Debug()).Msg("Teardown zdevice")
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *motionSensor) handleMessage(msg mqtt.Message) {
	var payload statusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	s.log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBatteryLevel(float32(payload.Battery)/float32(100), s)
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.Linkquality)/float32(255), s)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s)
	}

	if payload.Occupancy {
		s.systemControl.Alarm(system.BurglarAlarm, s)
	}
}
