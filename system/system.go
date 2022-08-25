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

type DeviceType string

const (
	Keypad        DeviceType = "keypad"
	MotionSensor  DeviceType = "motion"
	ContactSensor DeviceType = "contact"
	Siren         DeviceType = "siren"
	SmokeSensor   DeviceType = "smoke"
)

type Device struct {
	Id   string     `json:"id"`
	Name string     `json:"name"`
	Type DeviceType `json:"type"`
}

type StateUpdateFunc func(state State, armingMode ArmingMode, alarmType AlarmType)
type NotificationFunc func(note Notification) bool

type Controller interface {
	Arm(mode ArmingMode, dev Device) bool
	Disarm(pin string, dev Device) bool
	Alarm(aType AlarmType, dev Device) bool
	ReportBatteryLevel(level float32, dev Device)
	ReportLinkQuality(link float32, dev Device)

	GetState() State
	GetArmingMode() ArmingMode
	GetAlarmType() AlarmType

	SubscribeStateUpdate(id interface{}, fun StateUpdateFunc)
	UnsubscribeStateUpdate(id interface{})
}
