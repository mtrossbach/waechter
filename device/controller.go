package device

import "github.com/mtrossbach/waechter/system"

type SystemController interface {
	Arm(mode system.ArmingMode, dev system.Device) bool
	Disarm(pin string, dev system.Device) bool
	Alarm(aType system.AlarmType, dev system.Device) bool
	ReportBatteryLevel(level float32, dev system.Device)
	ReportLinkQuality(link float32, dev system.Device)

	GetState() system.State
	GetArmingMode() system.ArmingMode
	GetAlarmType() system.AlarmType

	SubscribeStateUpdate(id interface{}, fun system.StateUpdateFunc)
	UnsubscribeStateUpdate(id interface{})
}
