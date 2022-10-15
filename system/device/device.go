package device

import (
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system/zone"
	"github.com/rs/zerolog"
)

type Device struct {
	Id     Id             `json:"id"`
	Zone   zone.Id        `json:"zone"`
	Active bool           `json:"-"`
	Spec   Spec           `json:"-"`
	State  map[Sensor]any `json:"-"`
}

func DeviceFromConfig(config config.DeviceConfig) Device {
	return Device{
		Id:     Id(config.Id),
		Zone:   zone.Id(config.Zone),
		Active: false,
		Spec:   Spec{},
		State:  map[Sensor]any{},
	}
}

func DInfo(device *Device) *zerolog.Event {
	if device != nil {
		return appendDeviceInfo(*device, log.Info())
	}
	return log.Info()
}

func DDebug(device *Device) *zerolog.Event {
	if device != nil {
		return appendDeviceInfo(*device, log.Debug())
	}
	return log.Debug()
}

func DError(device *Device) *zerolog.Event {
	if device != nil {
		return appendDeviceInfo(*device, log.Error())
	}
	return log.Error()
}

func appendDeviceInfo(device Device, event *zerolog.Event) *zerolog.Event {
	var sensors []string
	for _, s := range device.Spec.Sensors {
		sensors = append(sensors, string(s))
	}
	return event.Str("_id", string(device.Id)).Str("_name", device.Spec.DisplayName).Str("_zone", string(device.Zone)).Strs("_sensors", sensors)
}
