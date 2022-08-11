package siren

import (
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/system"
)

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
	Warning warning `json:"warning"`
}

type warning struct {
	Mode            mode  `json:"mode"`
	Level           level `json:"level"`
	StrobeLevel     level `json:"strobe_level"`
	Strobe          bool  `json:"strobe"`
	StrobeDutyCycle int   `json:"strobe_duty_cycle"`
	Duration        int   `json:"duration"`
}

type statusPayload struct {
	Battery     int            `json:"battery"`
	Linkquality int            `json:"linkquality"`
	Warning     warningPayload `json:"warning"`
	Tamper      bool           `json:"tamper"`
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
		Warning: warning{
			Mode:            mode,
			Level:           level(config.GetString(cLevel)),
			StrobeLevel:     level(config.GetString(cLevel)),
			Strobe:          true,
			StrobeDutyCycle: 5,
			Duration:        5,
		},
	}
}
