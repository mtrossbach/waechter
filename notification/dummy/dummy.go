package dummy

/*
type dummy struct {
}

func New() *dummy {
	return &dummy{}
}

func (d *dummy) NotifyAlarm(recipient system.Recipient, systemName string, alarmType system.AlarmState, device *device.Device) bool {
	system.DInfo(device).Str("recipient", recipient.Name).Str("system", systemName).Str("alarm", string(alarmType)).Msg("Dummy alarm notification")
	return true
}

func (d *dummy) NotifyRecovery(recipient system.Recipient, systemName string, device *device.Device) bool {
	system.DInfo(device).Str("recipient", recipient.Name).Str("system", systemName).Msg("Dummy recovery notification")
	return true
}

func (d *dummy) NotifyLowBattery(recipient system.Recipient, systemName string, device *device.Device, batteryLevel float32) bool {
	system.DInfo(device).Str("recipient", recipient.Name).Str("system", systemName).Float32("battery", batteryLevel).Msg("Dummy battery notification")
	return true
}

func (d *dummy) NotifyLowLinkQuality(recipient system.Recipient, systemName string, device *device.Device, quality float32) bool {
	system.DInfo(device).Str("recipient", recipient.Name).Str("system", systemName).Float32("quality", quality).Msg("Dummy link quality notification")
	return true
}

func (d *dummy) NotifyAutoArm(recipient system.Recipient, systemName string) bool {
	log.Info().Str("recipient", recipient.Name).Str("system", systemName).Msg("Dummy auto-arm notification")
	return true
}

func (d *dummy) NotifyAutoDisarm(recipient system.Recipient, systemName string) bool {
	log.Info().Str("recipient", recipient.Name).Str("system", systemName).Msg("Dummy auto-disarm notification")
	return true
}
*/
