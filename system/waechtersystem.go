package system

import (
	"github.com/mtrossbach/waechter/misc"
)

type WaechterSystem struct {
	state      State
	subsystems []Subsystem
	devices    map[string]Device
}

func NewWaechterSystem() *WaechterSystem {
	return &WaechterSystem{
		state:      Disarmed,
		subsystems: []Subsystem{},
		devices:    make(map[string]Device),
	}
}

func (ws *WaechterSystem) RegisterSubsystem(subsystem Subsystem) {
	ws.subsystems = append(ws.subsystems, subsystem)
	subsystem.Start(ws)
}

func (ws *WaechterSystem) AddDevice(device Device) {
	ws.devices[device.GetId()] = device
	device.OnSystemStateChanged(ws.state)
	device.Setup(ws)
	misc.Log.Infof("Added device %v", DevDesc(device))
}

func (ws *WaechterSystem) RemoveDeviceById(id string) {
	dev, ok := ws.devices[id]
	if ok {
		dev.Teardown()
		delete(ws.devices, id)
	}
}

func (ws *WaechterSystem) HasDeviceId(id string) bool {
	_, ok := ws.devices[id]
	return ok
}

func (ws *WaechterSystem) GetDeviceIdsForSubsystem(subsystem string) []string {
	var devices []string
	for _, v := range ws.devices {
		if v.GetSubsystem() == subsystem {
			devices = append(devices, v.GetId())
		}
	}
	return devices
}

func (ws *WaechterSystem) GetDeviceById(id string) Device {
	return ws.devices[id]
}

func (ws *WaechterSystem) ReportBattery(device Device, battery float32) {
	misc.Log.Debugf("Got battery %v for %v", battery, DevDesc(device))
}

func (ws *WaechterSystem) ReportLinkQuality(device Device, linkquality float32) {
	misc.Log.Debugf("Got link quality %v for %v", linkquality, DevDesc(device))
}

func (ws *WaechterSystem) ReportTampered(device Device) {
	misc.Log.Debugf("Tamper alert %v", DevDesc(device))
}

func (ws *WaechterSystem) ReportTriggered(device Device) {
	misc.Log.Debugf("Trigger alert %v", DevDesc(device))
}

func (ws *WaechterSystem) RequestState(state State) {
	ws.setState(state)
}

func (ws *WaechterSystem) setState(state State) {
	misc.Log.Infof("State: %v", state)
	ws.state = state
	for _, device := range ws.devices {
		device.OnSystemStateChanged(ws.state)
	}
}

func (ws *WaechterSystem) GetState() State {
	return ws.state
}
