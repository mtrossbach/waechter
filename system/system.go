package system

type State string

const (
	DisarmedState   State = "Disarmed"
	ArmingState     State = "Arming"
	ArmedState      State = "Armed"
	EntryDelayState State = "EntryDelay"
)

type AlarmType string

const (
	NoAlarm      AlarmType = "NoAlarm"
	BurglarAlarm AlarmType = "BurglarAlarm"
	PanicAlarm   AlarmType = "PanicAlarm"
	FireAlarm    AlarmType = "FireAlarm"
	TamperAlarm  AlarmType = "TamperAlarm"
)

type ArmingMode string

const (
	StayMode ArmingMode = "Stay"
	AwayMode ArmingMode = "Away"
)

type DeviceType string

const (
	Keypad        DeviceType = "Keypad"
	MotionSensor  DeviceType = "MotionSensor"
	ContactSensor DeviceType = "ContactSensor"
	Siren         DeviceType = "Siren"
	SmokeSensor   DeviceType = "SmokeSensor"
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
