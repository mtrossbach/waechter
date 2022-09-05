package system

type Recipient struct {
	Name  string
	Phone string
	Lang  string
}

type NotificationAdapter interface {
	NotifyAlarm(recipient Recipient, systemName string, alarmType AlarmType, device *Device) bool
	NotifyRecovery(recipient Recipient, systemName string, device *Device) bool
	NotifyLowBattery(recipient Recipient, systemName string, device *Device, batteryLevel float32) bool
	NotifyLowLinkQuality(recipient Recipient, systemName string, device *Device, quality float32) bool
	NotifyAutoArm(recipient Recipient, systemName string) bool
	NotifyAutoDisarm(recipient Recipient, systemName string) bool
}
