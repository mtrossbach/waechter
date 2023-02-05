package system

import "github.com/mtrossbach/waechter/system/device"

const systemDeviceId device.Id = "_system_"

func systemDevice() *device.Device {
	return &device.Device{
		Id:     systemDeviceId,
		Zone:   "system",
		Active: true,
		Spec:   device.Spec{},
		State:  map[device.Sensor]any{},
	}
}
