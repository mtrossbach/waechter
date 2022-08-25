package zigbee2mqtt

import (
	"github.com/mtrossbach/waechter/internal/wslice"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/contactsensor"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/keypad"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/motionsensor"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/siren"
	"github.com/mtrossbach/waechter/system"
)

type ZDevice interface {
	OnDeviceAnnounced()
	UpdateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType)
	Setup(connector *connector.Connector, systemControl system.Controller)
	Teardown()
}

func createDevice(deviceInfo model.Z2MDeviceInfo) ZDevice {
	var exposes []string

	for _, e := range deviceInfo.Definition.Exposes {
		exposes = append(exposes, e.Property)
	}

	device := system.Device{
		Id:   ieee2Id(deviceInfo.IeeeAddress),
		Name: deviceInfo.FriendlyName,
		Type: "",
	}

	if wslice.ContainsAll(exposes, []string{"action_code", "action"}) {
		return keypad.New(device)
	} else if wslice.ContainsAll(exposes, []string{"warning"}) {
		return siren.New(device)
	} else if wslice.ContainsAll(exposes, []string{"contact"}) {
		return contactsensor.New(device)
	} else if wslice.ContainsAll(exposes, []string{"occupancy"}) {
		return motionsensor.New(device)
	} else {
		return nil
	}
}
