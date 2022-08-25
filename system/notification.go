package system

import "fmt"

type NotificationType string

const (
	AlarmNotification    NotificationType = "alarm"
	InfoNotification     NotificationType = "info"
	RecoveryNotification NotificationType = "recovery"
)

type Notification struct {
	Title       string
	Type        NotificationType
	Description string
}

func lowBatteryNotification(device Device, battery float32) *Notification {
	return &Notification{
		Title:       "Low Battery",
		Type:        InfoNotification,
		Description: fmt.Sprintf("The battery level of the device %v (%v) is at %.0f%%. Please check and replace the batteries.", device.Name, device.Id, battery*100),
	}
}

func lowLinkQualityNotification(device Device, link float32) *Notification {
	return &Notification{
		Title:       "Link Quality",
		Type:        InfoNotification,
		Description: fmt.Sprintf("The link quality of the device %v (%v) is at %.0f%%. Please check this device and try to reposition it.", device.Name, device.Id, link*100),
	}
}

func alarmNotification(alarmType AlarmType, dev Device) *Notification {
	if alarmType == NoAlarm {
		return nil
	}

	var notif Notification

	switch alarmType {
	case BurglarAlarm:
		notif = Notification{
			Title:       "Burglar Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) triggered an alarm!", dev.Name, dev.Id),
		}
	case FireAlarm:
		notif = Notification{
			Title:       "Fire Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) has detected fire/smoke!", dev.Name, dev.Id),
		}
	case PanicAlarm:
		notif = Notification{
			Title:       "Panic Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) triggered panic mode!", dev.Name, dev.Id),
		}
	case TamperAlarm:
		notif = Notification{
			Title:       "Tamper Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) is tampered!", dev.Name, dev.Id),
		}
	}

	return &notif
}

func recoveryNotification(device Device) *Notification {
	return &Notification{
		Title:       "Recovery",
		Type:        RecoveryNotification,
		Description: fmt.Sprintf("The system is defused by %v (%v).", device.Name, device.Id),
	}
}
