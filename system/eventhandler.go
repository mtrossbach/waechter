package system

import (
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
	"github.com/mtrossbach/waechter/system/device"
	"github.com/mtrossbach/waechter/system/zone"
	"time"
)

type EventHandler interface {
	SystemStarted()

	BatteryWarning(device device.Device, zone zone.Zone)
	LowBatteryLevel(device device.Device, zone zone.Zone, value byte)
	LowLinkQuality(device device.Device, zone zone.Zone, value byte)

	Alarm(device device.Device, zone zone.Zone, alarm alarm.Type)
	Recovered(device device.Device)
	EntryDelay(device device.Device, zone zone.Zone, delay time.Duration)

	Armed(device device.Device, mode arm.Mode, zones []zone.Zone, deviceIds []device.Id, autoArm bool, delay time.Duration)
	Disarmed(device device.Device, autoDisarm bool)

	AddedZone(zone zone.Zone)
	ChangedZone(zone zone.Zone)
	RemovedZone(zone zone.Zone)

	AddedDevice(device device.Device, zone zone.Zone)
	MovedDevice(device device.Device, oldZone zone.Zone, newZone zone.Zone)
	RemovedDevice(device device.Device, zone zone.Zone)
}
