package system

import (
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/device"
	"github.com/mtrossbach/waechter/system/zone"
)

type NotificationAdapter interface {
	Name() string
	NotifyAlarm(person config.Person, systemName string, a alarm.Type, device device.Spec, zone zone.Zone) bool
	NotifyRecovery(person config.Person, systemName string, device device.Spec, zone zone.Zone) bool
	NotifyLowBattery(person config.Person, systemName string, device device.Spec, zone zone.Zone, batteryLevel float32) bool
	NotifyLowLinkQuality(person config.Person, systemName string, device device.Spec, zone zone.Zone, quality float32) bool
	NotifyAutoArm(person config.Person, systemName string) bool
	NotifyAutoDisarm(person config.Person, systemName string) bool
}
