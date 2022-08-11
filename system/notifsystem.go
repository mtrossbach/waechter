package system

import (
	"fmt"
	"github.com/mtrossbach/waechter/misc"
	"github.com/rs/zerolog"
)

type notifSystem struct {
	notificationSubsystems []NotifSubsystem
	log                    zerolog.Logger
}

func newNotifSystem() *notifSystem {
	return &notifSystem{
		notificationSubsystems: []NotifSubsystem{},
		log:                    misc.Logger("NotifSystem"),
	}
}

func (ws *notifSystem) RegisterSubsystem(system NotifSubsystem) {
	ws.notificationSubsystems = append(ws.notificationSubsystems, system)
	ws.log.Info().Str("name", system.GetName()).Msg("Registered new external notif system")
}

func (ws *notifSystem) send(notification Notification) {
	ws.log.Info().Str("title", notification.Title).Str("type", string(notification.Type)).Str("description", notification.Description).Msg("Sending notifications ...")
	for _, ns := range ws.notificationSubsystems {
		ns.SendNotification(notification)
	}
}

func (ws *notifSystem) NotifyLowBattery(device Device, battery float32) {
	notification := Notification{
		Title:       "Low Battery",
		Type:        InfoNotification,
		Description: fmt.Sprintf("The battery level of the device %v (%v) is at %.0f%%. Please check and replace the batteries.", device.GetDisplayName(), device.GetId(), battery*100),
	}
	ws.send(notification)
}

func (ws *notifSystem) NotifyLowLinkQuality(device Device, link float32) {
	notification := Notification{
		Title:       "Link Quality",
		Type:        InfoNotification,
		Description: fmt.Sprintf("The link quality of the device %v (%v) is at %.0f%%. Please check this device and try to reposition it.", device.GetDisplayName(), device.GetId(), link*100),
	}
	ws.send(notification)
}

func (ws *notifSystem) NotifyAlarm(alarmType AlarmType, dev Device) {
	if alarmType == NoAlarm {
		return
	}

	var notif Notification

	switch alarmType {
	case BurglarAlarm:
		notif = Notification{
			Title:       "Burglar Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) triggered an alarm!", dev.GetDisplayName(), dev.GetId()),
		}
	case FireAlarm:
		notif = Notification{
			Title:       "Fire Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) has detected fire/smoke!", dev.GetDisplayName(), dev.GetId()),
		}
	case PanicAlarm:
		notif = Notification{
			Title:       "Panic Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) triggered panic mode!", dev.GetDisplayName(), dev.GetId()),
		}
	case TamperAlarm:
		notif = Notification{
			Title:       "Tamper Alert",
			Type:        AlarmNotification,
			Description: fmt.Sprintf("The device %v (%v) is tampered!", dev.GetDisplayName(), dev.GetId()),
		}
	}

	ws.send(notif)
}

func (ws *notifSystem) NotifyRecovery(device Device) {
	notification := Notification{
		Title:       "Recovery",
		Type:        RecoveryNotification,
		Description: fmt.Sprintf("The system is defused by %v (%v).", device.GetDisplayName(), device.GetId()),
	}
	ws.send(notification)
}
