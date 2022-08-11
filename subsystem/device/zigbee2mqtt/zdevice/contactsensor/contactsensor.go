package contactsensor

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type contactSensor struct {
	deviceInfo    model2.Z2MDeviceInfo
	z2mManager    *connector.Z2MManager
	systemControl system.Controller
	log           zerolog.Logger
}

func NewContactSensor(deviceInfo model2.Z2MDeviceInfo, z2mManager *connector.Z2MManager) *contactSensor {
	return &contactSensor{
		deviceInfo: deviceInfo,
		z2mManager: z2mManager,
		log:        misc.Logger("contactSensor"),
	}
}

func (s *contactSensor) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *contactSensor) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *contactSensor) GetSubsystem() string {
	return model2.SubsystemName
}

func (s *contactSensor) GetType() system.DeviceType {
	return system.ContactSensor
}

func (s *contactSensor) OnDeviceAnnounced() {

}

func (s *contactSensor) OnSystemStateChanged(state system.State, aMode system.ArmingMode, aType system.AlarmType) {

}

func (s *contactSensor) Setup(systemControl system.Controller) {
	system.DevLog(s, s.log.Debug()).Msg("Setup zdevice")
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *contactSensor) Teardown() {
	system.DevLog(s, s.log.Debug()).Msg("Teardown zdevice")
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *contactSensor) handleMessage(msg mqtt.Message) {
	var payload statusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		s.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	s.log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

	if payload.Battery > 0 {
		s.systemControl.ReportBatteryLevel(float32(payload.Battery)/float32(100), s)
	}

	if payload.Linkquality > 0 {
		s.systemControl.ReportLinkQuality(float32(payload.Linkquality)/float32(255), s)
	}

	if payload.Tamper {
		s.systemControl.Alarm(system.TamperAlarm, s)
	}

	if !payload.Contact {
		s.systemControl.Alarm(system.BurglarAlarm, s)
	}
}
