package zdevice

import (
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/contactsensor"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/keypad"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/motionsensor"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice/siren"
	"github.com/mtrossbach/waechter/system"
)

type ZDevice interface {
	system.Device
	OnDeviceAnnounced()
}

func CreateDevice(deviceInfo model.Z2MDeviceInfo, connector *connector.Connector) system.Device {
	var exposes []string

	for _, e := range deviceInfo.Definition.Exposes {
		exposes = append(exposes, e.Property)
	}

	if misc.ContainsAll(exposes, []string{"action_code", "action"}) {
		return keypad.New(deviceInfo, connector)
	} else if misc.ContainsAll(exposes, []string{"warning"}) {
		return siren.New(deviceInfo, connector)
	} else if misc.ContainsAll(exposes, []string{"contact"}) {
		return contactsensor.New(deviceInfo, connector)
	} else if misc.ContainsAll(exposes, []string{"occupancy"}) {
		return motionsensor.New(deviceInfo, connector)
	} else {
		return nil
	}
}
