package zigbee2mqtt

import (
	zdevice2 "github.com/mtrossbach/waechter/device/zigbee2mqtt/zdevice"
	"github.com/mtrossbach/waechter/internal/wslice"
	"github.com/mtrossbach/waechter/system"
)

func createDevice(deviceInfo Z2MDeviceInfo) ZDevice {
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
		return zdevice2.NewKeypad(device)
	} else if wslice.ContainsAll(exposes, []string{"warning"}) {
		return zdevice2.NewSiren(device)
	} else if wslice.ContainsAll(exposes, []string{"contact"}) {
		return zdevice2.NewContactSensor(device)
	} else if wslice.ContainsAll(exposes, []string{"occupancy"}) {
		return zdevice2.NewMotionSensor(device)
	} else {
		return nil
	}
}
