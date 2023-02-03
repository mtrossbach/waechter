package system

import (
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/device"
	"github.com/mtrossbach/waechter/system/zone"
)

type notificationManager struct {
	adapters []NotificationAdapter
}

func newNotificationManager() *notificationManager {
	return &notificationManager{adapters: []NotificationAdapter{}}
}

func (n *notificationManager) AddAdapter(adapter NotificationAdapter) {
	n.adapters = append(n.adapters, adapter)
}

func (n *notificationManager) allPersons() []config.Person {
	return config.Persons()
}

func (n *notificationManager) notify(persons []config.Person, handler func(person config.Person, adapter NotificationAdapter) bool) {
	for _, p := range persons {
		var successAdapter NotificationAdapter
		for _, a := range n.adapters {
			if handler(p, a) {
				successAdapter = a
				break
			}
		}
		if successAdapter != nil {
			log.Info().Str("name", p.Name).Str("adapter", successAdapter.Name()).Msg("Sent notification")
		}
	}
}

func (n *notificationManager) NotifyAlarm(a alarm.Type, device device.Spec, zone zone.Zone) {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyAlarm(person, config.SystemName(), a, device, zone)
	})
}

func (n *notificationManager) NotifyRecovery(device device.Spec, zone zone.Zone) {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyRecovery(person, config.SystemName(), device, zone)
	})
}

func (n *notificationManager) NotifyLowBattery(device device.Spec, zone zone.Zone, batteryLevel float32) {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyLowBattery(person, config.SystemName(), device, zone, batteryLevel)
	})
}

func (n *notificationManager) NotifyLowLinkQuality(device device.Spec, zone zone.Zone, quality float32) {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyLowLinkQuality(person, config.SystemName(), device, zone, quality)
	})
}

func (n *notificationManager) NotifyAutoArm() {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyAutoArm(person, config.SystemName())
	})
}

func (n *notificationManager) NotifyAutoDisarm() {
	n.notify(n.allPersons(), func(person config.Person, adapter NotificationAdapter) bool {
		return adapter.NotifyAutoDisarm(person, config.SystemName())
	})
}
