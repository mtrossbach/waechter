package system

import (
	"github.com/mtrossbach/waechter/system/device"
)

type DeviceConnector interface {
	Setup(controller Controller)
	Teardown()

	Id() string
	DisplayName() string

	Operational() bool

	EnumerateDevices() []device.Spec

	ActivateDevice(id device.Id) error
	DeactivateDevice(id device.Id) error

	ControlActor(id device.Id, actor device.Actor, value any) bool
}

type Controller interface {
	DeliverSensorValue(id device.Id, sensor device.Sensor, value any) bool

	DeviceListUpdated(connector DeviceConnector)

	OperationalStateChanged(connector DeviceConnector)

	DeviceUnavailable(id device.Id)
	DeviceAvailable(id device.Id)

	SystemState() State
}
