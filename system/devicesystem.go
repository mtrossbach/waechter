package system

import (
	"github.com/mtrossbach/waechter/misc"
	"github.com/rs/zerolog"
)

type deviceSystem struct {
	log        zerolog.Logger
	subsystems []DeviceSubsystem
	devices    map[string]Device
	controller Controller
}

func newDeviceSystem(controller Controller) *deviceSystem {
	return &deviceSystem{
		log:        misc.Logger("DeviceSystem"),
		subsystems: []DeviceSubsystem{},
		devices:    make(map[string]Device),
		controller: controller,
	}
}

func (ds *deviceSystem) RegisterSubsystem(subsystem DeviceSubsystem) {
	ds.subsystems = append(ds.subsystems, subsystem)
	ds.log.Info().Str("name", subsystem.GetName()).Msg("Registered new device subsystem")
	subsystem.Start(ds)
}

func (ds *deviceSystem) AddDevice(dev Device) {
	ds.devices[dev.GetId()] = dev
	ds.UpdateSystemStateOnDevice(dev)
	dev.Setup(ds.controller)
	DevLog(dev, ds.log.Info()).Msg("Added device")
}

func (ds *deviceSystem) RemoveDeviceById(id string) {
	dev, ok := ds.devices[id]
	if ok {
		dev.Teardown()
		delete(ds.devices, id)
	}
}

func (ds *deviceSystem) HasDeviceId(id string) bool {
	_, ok := ds.devices[id]
	return ok
}

func (ds *deviceSystem) GetDeviceIdsForSubsystem(name string) []string {
	var devices []string
	for _, v := range ds.devices {
		if v.GetSubsystem() == name {
			devices = append(devices, v.GetId())
		}
	}
	return devices
}

func (ds *deviceSystem) GetDeviceById(id string) Device {
	return ds.devices[id]
}

func (ds *deviceSystem) UpdateSystemState() {
	for _, dev := range ds.devices {
		ds.UpdateSystemStateOnDevice(dev)
	}
}

func (ds *deviceSystem) UpdateSystemStateOnDevice(dev Device) {
	if dev != nil {
		dev.OnSystemStateChanged(ds.controller.GetState(), ds.controller.GetArmingMode(), ds.controller.GetAlarmType())
	}
}
