package system

import (
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
	"github.com/mtrossbach/waechter/system/device"
	"time"
)

type State struct {
	ArmMode arm.Mode
	Alarm   alarm.Type

	armModeUpdated time.Time
}

func (s State) Armed() bool {
	return s.ArmMode != arm.Disarmed
}

func (s State) stateActorPayload() device.StateActorPayload {
	return device.StateActorPayload{
		ArmMode: s.ArmMode,
		Alarm:   s.Alarm,
	}
}

func (s State) alarmActorPayload() device.AlarmActorPayload {
	return device.AlarmActorPayload{Alarm: s.Alarm}
}
