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

type Device interface {
	GetId() string
	GetDisplayName() string
	GetSubsystem() string
	GetType() DeviceType

	OnSystemStateChanged(state State, aMode ArmingMode, aType AlarmType)
	Setup(Controller)
	Teardown()
}

type Controller interface {
	Arm(mode ArmingMode, dev Device) bool
	Disarm(pin string, dev Device) bool
	Alarm(aType AlarmType, dev Device) bool
	ReportBatteryLevel(level float32, dev Device)
	ReportLinkQuality(link float32, dev Device)

	GetState() State
	GetArmingMode() ArmingMode
	GetAlarmType() AlarmType
}

type DeviceSystem interface {
	RegisterSubsystem(subsystem DeviceSubsystem)
	AddDevice(dev Device)
	RemoveDeviceById(id string)
	HasDeviceId(id string) bool
	GetDeviceIdsForSubsystem(name string) []string
	GetDeviceById(id string) Device
}

type DeviceSubsystem interface {
	GetName() string
	Start(DeviceSystem)
	Stop()
}

type NotificationType string

const (
	AlarmNotification    NotificationType = "alarm"
	InfoNotification     NotificationType = "info"
	RecoveryNotification NotificationType = "recovery"
)

type Notification struct {
	Title       string
	Type        NotificationType
	Description string
}

type NotifSubsystem interface {
	GetName() string
	SendNotification(Notification)
}
