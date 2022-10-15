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
