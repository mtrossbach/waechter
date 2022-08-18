package siren

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/cfg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type siren struct {
	deviceInfo    model2.Z2MDeviceInfo
	connector     *connector.Connector
	systemControl system.Controller
	log           zerolog.Logger
	targetTopic   string
}

func New(deviceInfo model2.Z2MDeviceInfo, connector *connector.Connector) *siren {
	return &siren{
		deviceInfo:  deviceInfo,
		connector:   connector,
		log:         cfg.Logger("Z2MSiren"),
		targetTopic: fmt.Sprintf("%v/set", deviceInfo.FriendlyName),
	}
}

func (s *siren) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *siren) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *siren) GetSubsystem() string {
	return model2.SubsystemName
}

func (s *siren) GetType() system.DeviceType {
	return system.Siren
}

func (s *siren) OnSystemStateChanged(state system.State, aMode system.ArmingMode, aType system.AlarmType) {
	s.sendState()
}

func (s *siren) OnDeviceAnnounced() {
	s.sendState()
}

func (s *siren) sendState() {
	var payload warningPayload
	if s.systemControl.GetAlarmType() != system.NoAlarm && cfg.GetBool(cEnabled) {
		payload = newWarningPayload(s.systemControl.GetAlarmType())
	} else {
		payload = newWarningPayload(system.NoAlarm)
	}
	s.connector.Publish(s.targetTopic, payload)
}

func (s *siren) Setup(systemControl system.Controller) {
	system.DevLog(s, s.log.Debug()).Msg("Setup zdevice")
	s.systemControl = systemControl
	s.connector.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *siren) Teardown() {
	system.DevLog(s, s.log.Debug()).Msg("Teardown zdevice")
	s.systemControl = nil
	s.connector.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *siren) handleMessage(msg mqtt.Message) {
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
}
