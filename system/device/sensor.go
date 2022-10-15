package device

import (
	"github.com/mtrossbach/waechter/system/arm"
)

type Sensor string

const (
	MotionSensor         Sensor = "motion"
	ContactSensor        Sensor = "contact"
	SmokeSensor          Sensor = "smoke"
	PanicSensor          Sensor = "panic"
	BatteryWarningSensor Sensor = "battery-warning"
	TamperSensor         Sensor = "tamper"
	BatteryLevelSensor   Sensor = "battery-level"
	LinkQualitySensor    Sensor = "link-quality"
	ArmingSensor         Sensor = "arming"
	DisarmingSensor      Sensor = "disarming"
)

type MotionSensorValue struct {
	Motion bool
}

type ContactSensorValue struct {
	Contact bool
}

type SmokeSensorValue struct {
	Smoke bool
}

type PanicSensorValue struct {
	Panic bool
}

type BatteryWarningSensorValue struct {
	BatteryWarning bool
}

type TamperSensorValues struct {
	Tamper bool
}

type BatteryLevelSensorValue struct {
	BatteryLevel float32
}

type LinkQualitySensorValue struct {
	LinkQuality float32
}

type ArmingSensorValue struct {
	ArmMode arm.Mode
}

type DisarmingSensorValue struct {
	Pin string
}
