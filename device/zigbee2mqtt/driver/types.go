package driver

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/system"
)

type Sender func(payload any)

type mode string
type level string

const (
	stop      mode = "stop"
	burglar   mode = "burglar"
	fire      mode = "fire"
	emergency mode = "emergency"
)

const (
	low      level = "low"
	medium   level = "medium"
	high     level = "high"
	veryHigh level = "very_high"
)

type baseStatus struct {
	Battery     int  `json:"battery"`
	LinkQuality int  `json:"linkquality"`
	Tamper      bool `json:"tamper"`
}

type sirenStatus struct {
	baseStatus
	Warning warningPayload `json:"warning"`
}

type smokeStatus struct {
	baseStatus
	Smoke bool `json:"smoke"`
}

type motionStatus struct {
	baseStatus
	Occupancy bool `json:"occupancy"`
}

type contactStatus struct {
	baseStatus
	Contact bool `json:"contact"`
}

type keypadStatus struct {
	motionStatus
	Action            string `json:"action"`
	ActionCode        string `json:"action_code"`
	ActionTransaction int    `json:"action_transaction"`
	ActionZone        int    `json:"action_zone"`
}

type keypadSetState struct {
	ArmMode armMode `json:"arm_mode"`
}
type armMode struct {
	Mode        string `json:"mode"`
	Transaction *int   `json:"transaction,omitempty"`
}

type warningPayload struct {
	Warning warningOptions `json:"warning"`
}

type warningOptions struct {
	Mode            mode  `json:"mode"`
	Level           level `json:"level"`
	StrobeLevel     level `json:"strobe_level"`
	Strobe          bool  `json:"strobe"`
	StrobeDutyCycle int   `json:"strobe_duty_cycle"`
	Duration        int   `json:"duration"`
}

func newWarningPayload(alarmType system.AlarmType) warningPayload {
	mode := stop
	switch alarmType {
	case system.BurglarAlarm, system.TamperAlarm:
		mode = burglar
	case system.PanicAlarm:
		mode = emergency
	case system.FireAlarm:
		mode = fire
	default:
		mode = stop
	}

	return warningPayload{
		Warning: warningOptions{
			Mode:            mode,
			Level:           level(cfg.GetString(cLevel)),
			StrobeLevel:     level(cfg.GetString(cLevel)),
			Strobe:          mode != stop,
			StrobeDutyCycle: 1,
			Duration:        10 * 60,
		},
	}
}

func newKeypadSetState(mode string, transaction *int) keypadSetState {
	return keypadSetState{ArmMode: armMode{Mode: mode, Transaction: transaction}}
}

func sysStateToKeypadState(state system.State, aMode system.ArmingMode, aType system.AlarmType) string {
	switch aType {
	case system.NoAlarm:
		break
	case system.PanicAlarm:
		return "panic"
	default:
		return "in_alarm"
	}

	switch state {
	case system.DisarmedState:
		return "disarm"
	case system.ArmingState:
		return "exit_delay"
	case system.ArmedState:
		switch aMode {
		case system.AwayMode:
			return "arm_all_zones"
		case system.StayMode:
			return "arm_day_zones"
		}
	case system.EntryDelayState:
		return "entry_delay"
	}
	return "disarm"
}
