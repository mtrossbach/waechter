package device

import (
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
)

type Actor string

const (
	AlarmActor             Actor = "alarm"
	StateActor             Actor = "state"
	NotificationShortActor Actor = "notification-short"
	NotificationLongActor  Actor = "notification-long"
)

type AlarmActorPayload struct {
	Alarm alarm.Type
}

type StateActorPayload struct {
	ArmMode arm.Mode
	Alarm   alarm.Type
}
