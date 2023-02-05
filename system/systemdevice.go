package system

import "github.com/mtrossbach/waechter/system/device"

const systemDeviceId device.Id = "_system_"

func systemDevice() *device.Device {
	return &device.Device{
		Id:     systemDeviceId,
		Zone:   "system",
		Active: true,
		Spec: device.Spec{
			Id:          systemDeviceId,
			DisplayName: "System Device",
			Vendor:      "",
			Model:       "",
			Description: "",
			Sensors:     nil,
			Actors:      nil,
		},
		State: map[device.Sensor]any{},
	}
}
