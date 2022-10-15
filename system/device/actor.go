package device

import (
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
)

type Actor string

const (
	AlarmActor Actor = "alarm"
	StateActor Actor = "state"
)

type AlarmActorPayload struct {
	Alarm alarm.Type
}

type StateActorPayload struct {
	ArmMode arm.Mode
	Alarm   alarm.Type
}
