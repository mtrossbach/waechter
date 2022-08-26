package system

type State string

const (
	DisarmedState   State = "disarmed"
	ArmingState     State = "arming"
	ArmedState      State = "armed"
	EntryDelayState State = "entry-delay"
)

type AlarmType string

const (
	NoAlarm      AlarmType = "none"
	BurglarAlarm AlarmType = "burglar"
	PanicAlarm   AlarmType = "panic"
	FireAlarm    AlarmType = "fire"
	TamperAlarm  AlarmType = "tamper"
)

type ArmingMode string

const (
	StayMode ArmingMode = "stay"
	AwayMode ArmingMode = "away"
)

type StateUpdateFunc func(state State, armingMode ArmingMode, alarmType AlarmType)
type NotificationFunc func(note Notification) bool
