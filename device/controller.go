package device

import "github.com/mtrossbach/waechter/system"

type SystemController interface {
	ArmStay(dev system.Device) bool
	ArmAway(dev system.Device) bool
	Disarm(pin string, dev system.Device) bool
	ForceDisarm(dev system.Device)
	Alarm(aType system.AlarmType, dev system.Device) bool
	ReportBatteryLevel(level float32, dev system.Device)
	ReportLinkQuality(link float32, dev system.Device)

	GetArmState() system.ArmState
	GetAlarmType() system.AlarmType

	SubscribeStateUpdate(id interface{}, fun system.StateUpdateFunc)
	UnsubscribeStateUpdate(id interface{})
}
