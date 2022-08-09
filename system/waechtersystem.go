package system

import (
	"time"

	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
	"github.com/rs/zerolog"
)

type WaechterSystem struct {
	state      State
	subsystems []Subsystem
	devices    map[string]Device
	log        zerolog.Logger
}

func NewWaechterSystem() *WaechterSystem {
	return &WaechterSystem{
		state:      Disarmed,
		subsystems: []Subsystem{},
		devices:    make(map[string]Device),
		log:        misc.Logger("WaechterSystem"),
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
	ws.log.Info().Str("device", DevDesc(device)).Msg("Added device.")
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
	ws.log.Debug().Float32("battery", battery).Str("device", DevDesc(device)).Msg("Got battery info.")
}

func (ws *WaechterSystem) ReportLinkQuality(device Device, linkquality float32) {
	ws.log.Debug().Float32("link", linkquality).Str("device", DevDesc(device)).Msg("Got link quality info.")
}

func (ws *WaechterSystem) ReportTampered(device Device) {
	ws.log.Debug().Str("device", DevDesc(device)).Msg("Tamper alert!")
}

func (ws *WaechterSystem) ReportTriggered(device Device) {
	ws.log.Debug().Str("device", DevDesc(device)).Msg("Triggered!")
}

func (ws *WaechterSystem) setState(state State) {
	ws.log.Info().Str("state", string(state)).Msg("Updated state")
	ws.state = state
	ws.notifyState()
}

func (ws *WaechterSystem) notifyState() {
	for _, device := range ws.devices {
		device.OnSystemStateChanged(ws.state)
	}
}

func (ws *WaechterSystem) GetState() State {
	return ws.state
}

func (ws *WaechterSystem) ArmStay() {
	if ws.state == Disarmed {
		ws.setState(ArmingStay)
		ws.armingTimer()
	} else {
		ws.notifyState()
	}
}

func (ws *WaechterSystem) ArmAway() {
	if ws.state == Disarmed {
		ws.setState(ArmingAway)
		ws.armingTimer()
	} else {
		ws.notifyState()
	}
}

func (ws *WaechterSystem) Disarm(enteredPin string) {
	pins := config.GetConfig().DisarmPins

	for _, pin := range pins {
		if pin.Pin == enteredPin {
			ws.setState(Disarmed)
			return
		}
	}

	ws.notifyState()
}

func (ws *WaechterSystem) Panic() {
	ws.setState(Panic)
}

func (ws *WaechterSystem) armingTimer() {
	time.AfterFunc(config.GetConfig().General.ExitDelay, func() {
		if ws.state == ArmingAway {
			ws.setState(ArmedAway)
		} else if ws.state == ArmingStay {
			ws.setState(ArmedStay)
		}
	})
}
