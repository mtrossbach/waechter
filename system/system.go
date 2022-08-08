package system

import "fmt"

type State string

//               ArmingStay -> ArmedStay -> EntryDelay -> InAlarm ->
//             /              \            \             \           \
// Disarmed ->                 ->------------>------------->----------->-------> Disarmed
//             \              /            /             /           /
//               ArmingAway -> ArmedAway -> EntryDelay -> InAlarm ->
const (
	Disarmed State = "Disarmed"

	ArmingStay State = "ArmingStay"
	ArmingAway State = "ArmingAway"

	ArmedStay State = "ArmedStay"
	ArmedAway State = "ArmedAway"

	InAlarm State = "InAlarm"

	EntryDelay State = "EntryDelay"
)

type DeviceType string

const (
	Unknown       DeviceType = ""
	Keypad        DeviceType = "keypad"
	MotionSensor  DeviceType = "motionsensor"
	ContactSensor DeviceType = "contactsensor"
	Siren         DeviceType = "siren"

	/*
		GasSensor
		SmokeSensor
		WaterLeakSensor
	*/
)

type Device interface {
	GetId() string
	GetDisplayName() string
	GetSubsystem() string
	GetType() DeviceType

	OnSystemStateChanged(State)
	Setup(SystemControl)
	Teardown()
}

func DevDesc(device Device) string {
	return fmt.Sprintf("[%v:%v:%v:%v]", device.GetType(), device.GetId(), device.GetDisplayName(), device.GetSubsystem())
}

type SystemControl interface {
	ReportBattery(Device, float32)
	ReportLinkQuality(Device, float32)
	ReportTampered(Device)
	ReportTriggered(Device)

	RequestState(State)
	GetState() State
}

type DeviceManager interface {
	AddDevice(Device)
	RemoveDeviceById(string)
	HasDeviceId(string) bool
	GetDeviceIdsForSubsystem(string) []string
	GetDeviceById(string) Device
}

type Subsystem interface {
	GetName() string
	Start(DeviceManager)
	Stop()
}
