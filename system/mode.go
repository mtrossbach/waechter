package system

type State string

const (
	DisarmedState   State = "disarmed"
	ArmingState     State = "arming"
	ArmedState      State = "armed"
	EntryDelayState State = "entry-delay"
)

func isValidState(s State) bool {
	return s == DisarmedState || s == ArmingState || s == ArmedState || s == EntryDelayState
}

type AlarmType string

const (
	NoAlarm      AlarmType = "none"
	BurglarAlarm AlarmType = "burglar"
	PanicAlarm   AlarmType = "panic"
	FireAlarm    AlarmType = "fire"
	TamperAlarm  AlarmType = "tamper"
)

func isValidAlarmType(a AlarmType) bool {
	return a == NoAlarm || a == BurglarAlarm || a == PanicAlarm || a == FireAlarm || a == TamperAlarm
}

type ArmingMode string

const (
	StayMode ArmingMode = "stay"
	AwayMode ArmingMode = "away"
)

func isValidArmingMode(a ArmingMode) bool {
	return a == StayMode || a == AwayMode
}
