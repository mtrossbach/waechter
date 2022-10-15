package zigbee2mqtt

import (
	"github.com/mtrossbach/waechter/system"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
)

func extract[T any](m map[string]any, key string) *T {
	a, ok := m[key]
	if ok {
		if v, ok := a.(T); ok {
			return &v
		}
	}
	return nil
}

type armModePayload struct {
	ArmMode armMode `json:"arm_mode"`
}

type armMode struct {
	Mode        string `json:"mode"`
	Transaction *int   `json:"transaction,omitempty"`
}

func newArmModePayload(state system.State, transactionId *int) armModePayload {
	var mode string

	switch state.Alarm {
	case alarm.None:
		switch state.ArmMode {
		case arm.Disarmed:
			mode = "disarm"
		case arm.Perimeter:
			mode = "arm_day_zones"
		default:
			mode = "arm_all_zones"
		}
	case alarm.EntryDelay:
		mode = "entry_delay"
	case alarm.Panic:
		mode = "panic"
	default:
		mode = "in_alarm"
	}

	return armModePayload{ArmMode: armMode{
		Mode:        mode,
		Transaction: transactionId,
	}}
}

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

func newWarningPayload(a alarm.Type) warningPayload {
	mode := stop
	switch a {
	case alarm.Burglar, alarm.Tamper, alarm.TamperPin:
		mode = burglar
	case alarm.Panic:
		mode = emergency
	case alarm.Fire:
		mode = fire
	default:
		mode = stop
	}

	return warningPayload{
		Warning: warningOptions{
			Mode:            mode,
			Level:           high,
			StrobeLevel:     high,
			Strobe:          mode != stop,
			StrobeDutyCycle: 1,
			Duration:        10 * 60,
		},
	}
}
