package system

type ArmState string

const (
	DisarmedState  ArmState = "disarmed"
	ExitDelayState ArmState = "exit-delay"
	ArmedStayState ArmState = "armed-stay"
	ArmedAwayState ArmState = "armed-away"
)

func isValidArmState(s ArmState) bool {
	return s == DisarmedState || s == ExitDelayState || s == ArmedStayState || s == ArmedAwayState
}

type AlarmType string

const (
	NoAlarm         AlarmType = "none"
	BurglarAlarm    AlarmType = "burglar"
	PanicAlarm      AlarmType = "panic"
	FireAlarm       AlarmType = "fire"
	TamperAlarm     AlarmType = "tamper"
	EntryDelayAlarm AlarmType = "entry-delay"
)

func isValidAlarmType(a AlarmType) bool {
	return a == NoAlarm || a == BurglarAlarm || a == PanicAlarm || a == FireAlarm || a == TamperAlarm || a == EntryDelayAlarm
}
