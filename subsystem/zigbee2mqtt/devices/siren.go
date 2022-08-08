package devices

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
)

type genericSiren struct {
	deviceInfo    model.Z2MDeviceInfo
	z2mManager    *zigbee.Z2MManager
	systemControl system.SystemControl
}

func newGenericSiren(deviceInfo model.Z2MDeviceInfo, z2mManager *zigbee.Z2MManager) *genericSiren {
	return &genericSiren{
		deviceInfo: deviceInfo,
		z2mManager: z2mManager,
	}
}

func (s *genericSiren) GetId() string {
	return s.deviceInfo.IeeeAddress
}

func (s *genericSiren) GetDisplayName() string {
	return s.deviceInfo.FriendlyName
}

func (s *genericSiren) GetSubsystem() string {
	return model.SubsystemName
}

func (s *genericSiren) GetType() system.DeviceType {
	return system.Siren
}

func (s *genericSiren) OnSystemStateChanged(state system.State) {

}

func (s *genericSiren) OnDeviceAnnounced() {

}

func (s *genericSiren) Setup(systemControl system.SystemControl) {
	misc.Log.Debugf("Setup device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = systemControl
	s.z2mManager.Subscribe(s.deviceInfo.FriendlyName, s.handleMessage)
}

func (s *genericSiren) Teardown() {
	misc.Log.Debugf("Teardown device %v:%v:%v", s.GetType(), s.GetId(), s.GetDisplayName())
	s.systemControl = nil
	s.z2mManager.Unsubscribe(s.deviceInfo.FriendlyName)
}

func (s *genericSiren) handleMessage(msg mqtt.Message) {
	misc.Log.Debugf("Got data: %v", string(msg.Payload()))
	//TODO
}
