package keypad

import "github.com/mtrossbach/waechter/system"

type statusPayload struct {
	Action            string `json:"action"`
	ActionCode        string `json:"action_code"`
	ActionTransaction int    `json:"action_transaction"`
	ActionZone        int    `json:"action_zone"`
	Battery           int    `json:"battery"`
	Linkquality       int    `json:"linkquality"`
	Occupancy         bool   `json:"occupancy"`
	Tamper            bool   `json:"tamper"`
}

type statePayload struct {
	ArmMode armMode `json:"arm_mode"`
}
type armMode struct {
	Mode        string `json:"mode"`
	Transaction *int   `json:"transaction,omitempty"`
}

func newStatePayload(mode string, transaction *int) statePayload {
	return statePayload{ArmMode: armMode{Mode: mode, Transaction: transaction}}
}

func systemStateToDeviceState(state system.State, aMode system.ArmingMode, aType system.AlarmType) string {
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
