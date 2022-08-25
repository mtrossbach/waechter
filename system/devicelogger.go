package system

import (
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/rs/zerolog"
)

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
	return event.Str("device", fmt.Sprintf("[%v;%v;%v]", device.Id, device.Name, device.Type))
}
