package contactsensor

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/system"
)

type contactSensor struct {
	deviceInfo    model2.Z2MDeviceInfo
	connector     *connector.Connector
	systemControl system.Controller
}

func New(deviceInfo model2.Z2MDeviceInfo, connector *connector.Connector) *contactSensor {
	return &contactSensor{
		deviceInfo: deviceInfo,
		connector:  connector,
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
	s.systemControl = systemControl
	s.connector.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *contactSensor) Teardown() {
	s.systemControl = nil
	s.connector.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *contactSensor) handleMessage(msg mqtt.Message) {
	var payload statusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
		return
	}

	log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

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
