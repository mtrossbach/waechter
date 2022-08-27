package system

import (
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/rs/zerolog"
)

type DeviceType string

const (
	Keypad        DeviceType = "keypad"
	MotionSensor  DeviceType = "motion"
	ContactSensor DeviceType = "contact"
	Siren         DeviceType = "siren"
	SmokeSensor   DeviceType = "smoke"
)

type Device struct {
	Namespace string     `json:"namespace"`
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Type      DeviceType `json:"type"`
}

func DInfo(device Device) *zerolog.Event {
	return appendDeviceInfo(device, log.Info())
}

func DDebug(device Device) *zerolog.Event {
	return appendDeviceInfo(device, log.Debug())
}

func DError(device Device) *zerolog.Event {
	return appendDeviceInfo(device, log.Error())
}

func appendDeviceInfo(device Device, event *zerolog.Event) *zerolog.Event {
	return event.Str("id", fmt.Sprintf("%v::%v", device.Namespace, device.Id)).Str("name", device.Name).Str("type", string(device.Type))
}
