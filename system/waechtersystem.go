package system

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"time"

	"github.com/rs/zerolog"
)

type WaechterSystem struct {
	notifSystem   notifSystem
	deviceSystem  deviceSystem
	state         State
	armingMode    ArmingMode
	alarmType     AlarmType
	log           zerolog.Logger
	wrongPinCount int
}

func NewWaechterSystem() *WaechterSystem {
	system := &WaechterSystem{
		state:         DisarmedState,
		armingMode:    AwayMode,
		alarmType:     NoAlarm,
		log:           cfg.Logger("WaechterSystem"),
		wrongPinCount: 0,
	}

	system.deviceSystem = *newDeviceSystem(system)
	system.notifSystem = *newNotifSystem()
	return system
}

func (ws *WaechterSystem) RegisterDeviceSubsystem(subsystem DeviceSubsystem) {
	ws.deviceSystem.RegisterSubsystem(subsystem)
}

func (ws *WaechterSystem) RegisterNotifSubsystem(subsystem NotifSubsystem) {
	ws.notifSystem.RegisterSubsystem(subsystem)
}

func (ws *WaechterSystem) Arm(mode ArmingMode, dev Device) bool {
	if ws.state == DisarmedState {
		ws.setState(ArmingState, mode, NoAlarm)
		time.AfterFunc(time.Duration(cfg.GetInt(cExitDelay))*time.Second, func() {
			if ws.state == ArmingState {
				ws.setState(ArmedState, ws.armingMode, NoAlarm)
			}
		})
		return true
	} else {
		// Requesting device probably has a wrong system state -> update it
		ws.deviceSystem.UpdateSystemStateOnDevice(dev)
		return false
	}
}

func (ws *WaechterSystem) Disarm(enteredPin string, dev Device) bool {
	if ws.state == DisarmedState {
		// Requesting device probably has a wrong system state -> update it
		ws.deviceSystem.UpdateSystemStateOnDevice(dev)
		return false
	}

	pinOk := false
	disarmPins := cfg.GetStrings(cDisarmPins)
	for _, pin := range disarmPins {
		if pin == enteredPin {
			pinOk = true
			break
		}
	}

	if pinOk {
		ws.wrongPinCount = 0
		if ws.alarmType != NoAlarm {
			ws.notifSystem.NotifyRecovery(dev)
		}
		ws.setState(DisarmedState, ws.armingMode, NoAlarm)
		return true
	} else {
		ws.wrongPinCount += 1
		if ws.wrongPinCount > cfg.GetInt(cMaxWrongPinCount) {
			ws.Alarm(BurglarAlarm, dev)
		}
		return false
	}
}

func (ws *WaechterSystem) Alarm(aType AlarmType, dev Device) bool {
	if (ws.state != ArmedState) && aType == BurglarAlarm {
		return false
	}
	if aType == TamperAlarm && !cfg.GetBool(cTamperAlarm) {
		return false
	}

	if ws.state == ArmedState && aType == BurglarAlarm {
		ws.setState(EntryDelayState, ws.armingMode, ws.alarmType)
		time.AfterFunc(time.Duration(cfg.GetInt(cEntryDelay))*time.Second, func() {
			if ws.state == EntryDelayState {
				ws.setState(ws.state, ws.armingMode, aType)
				ws.notifSystem.NotifyAlarm(aType, dev)
			}
		})
		return true
	} else {
		ws.setState(ws.state, ws.armingMode, aType)
		ws.notifSystem.NotifyAlarm(aType, dev)
		return true
	}
}

func (ws *WaechterSystem) ReportBatteryLevel(level float32, dev Device) {
	if level < cfg.GetFloat32(cBatteryThreshold) {
		DevLog(dev, ws.log.Info()).Float32("battery", level).Msg("Battery is too low! Notify!")
		ws.notifSystem.NotifyLowBattery(dev, level)
	} else {
		DevLog(dev, ws.log.Debug()).Float32("battery", level).Msg("Got battery info")
	}
}

func (ws *WaechterSystem) ReportLinkQuality(link float32, dev Device) {
	if ws.IsArmed() && link < cfg.GetFloat32(cLinkQualityThreshold) && cfg.GetBool(cTamperAlarm) {
		DevLog(dev, ws.log.Info()).Float32("link", link).Msg("Link quality is too low! Tamper alarm!")
		ws.Alarm(TamperAlarm, dev)
	} else if link < cfg.GetFloat32(cLinkQualityThreshold) {
		DevLog(dev, ws.log.Info()).Float32("link", link).Msg("Link quality is too low! Notify!")
		ws.notifSystem.NotifyLowLinkQuality(dev, link)
	} else {
		DevLog(dev, ws.log.Debug()).Float32("link", link).Msg("Got link quality info")
	}
}

func (ws *WaechterSystem) GetState() State {
	return ws.state
}

func (ws *WaechterSystem) GetArmingMode() ArmingMode {
	return ws.armingMode
}

func (ws *WaechterSystem) GetAlarmType() AlarmType {
	return ws.alarmType
}

func (ws *WaechterSystem) setState(state State, mode ArmingMode, alarmType AlarmType) {
	ws.state = state
	ws.armingMode = mode
	ws.alarmType = alarmType
	ws.log.Info().Str("state", string(state)).Str("armingMode", string(mode)).Str("alarmType", string(alarmType)).Msg("System state updated")
	ws.deviceSystem.UpdateSystemState()
}

func (ws *WaechterSystem) IsArmed() bool {
	return ws.state == ArmedState || ws.state == EntryDelayState
}
