package system

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
)

type notificationManager struct {
	systemName string
	adapters   []NotificationAdapter
}

func newNotificationManager(systemName string) *notificationManager {
	return &notificationManager{systemName: systemName, adapters: []NotificationAdapter{}}
}

func (n *notificationManager) AddAdapter(adapter NotificationAdapter) {
	n.adapters = append(n.adapters, adapter)
}

func (n *notificationManager) allRecipients() []Recipient {
	objs := cfg.GetStringStringMaps(cNotificationRecipients)
	var recipients []Recipient
	for _, m := range objs {
		recipients = append(recipients, Recipient{
			Name:  m[cRecipientName],
			Phone: m[cRecipientPhone],
			Lang:  m[cRecipientLanguage],
		})
	}
	return recipients
}

func (n *notificationManager) notify(recipients []Recipient, handler func(recipient Recipient, adapter NotificationAdapter) bool) {
	for _, r := range recipients {
		success := false
		for _, a := range n.adapters {
			if handler(r, a) {
				success = true
				break
			}
		}
		if !success {
			log.Error().Str("name", r.Name).Msg("Could not sent notification to this recipient, because all notification methods failed.")
		}
	}
}

func (n *notificationManager) notifyAlarm(alarmType AlarmType, device *Device) {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyAlarm(recipient, n.systemName, alarmType, device)
	})
}

func (n *notificationManager) notifyRecovery(device *Device) {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyRecovery(recipient, n.systemName, device)
	})
}

func (n *notificationManager) notifyLowBattery(device *Device, battery float32) {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyLowBattery(recipient, n.systemName, device, battery)
	})
}

func (n *notificationManager) notifyLowLinkQuality(device *Device, link float32) {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyLowLinkQuality(recipient, n.systemName, device, link)
	})
}

func (n *notificationManager) notifyAutoArmed() {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyAutoArm(recipient, n.systemName)
	})
}

func (n *notificationManager) notifyAutoDisarmed() {
	n.notify(n.allRecipients(), func(recipient Recipient, adapter NotificationAdapter) bool {
		return adapter.NotifyAutoDisarm(recipient, n.systemName)
	})
}
