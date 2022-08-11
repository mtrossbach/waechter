package system

import (
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
	"github.com/rs/zerolog"
	"time"
)

type WaechterSystem struct {
	notifSystem   notifSystem
	deviceSystem  deviceSystem
	config        config.General
	pins          []config.DisarmPin
	state         State
	armingMode    ArmingMode
	alarmType     AlarmType
	log           zerolog.Logger
	wrongPinCount int
}

func NewWaechterSystem() *WaechterSystem {
	system := &WaechterSystem{
		config:        config.GetConfig().General,
		pins:          config.GetConfig().DisarmPins,
		state:         DisarmedState,
		armingMode:    AwayMode,
		alarmType:     NoAlarm,
		log:           misc.Logger("WaechterSystem"),
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
		time.AfterFunc(ws.config.ExitDelay, func() {
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
	for _, pin := range ws.pins {
		if pin.Pin == enteredPin {
			pinOk = true
			break
		}
	}

	if pinOk {
		ws.wrongPinCount = 0
		if ws.state == InAlarmState {
			ws.notifSystem.NotifyRecovery(dev)
		}
		ws.setState(DisarmedState, ws.armingMode, NoAlarm)
		return true
	} else {
		ws.wrongPinCount += 1
		if ws.wrongPinCount > ws.config.MaxWrongPinCount {
			ws.Alarm(TamperAlarm, dev)
		}
		return false
	}
}

func (ws *WaechterSystem) Alarm(aType AlarmType, dev Device) bool {
	if (ws.state != ArmedState && ws.state != InAlarmState) && aType == BurglarAlarm {
		return false
	}
	if aType == TamperAlarm && !ws.config.TamperAlarm {
		return false
	}

	if ws.state == ArmedState && aType == BurglarAlarm {
		ws.setState(EntryDelayState, ws.armingMode, aType)
		time.AfterFunc(ws.config.EntryDelay, func() {
			if ws.state == EntryDelayState {
				ws.setState(InAlarmState, ws.armingMode, aType)
				ws.notifSystem.NotifyAlarm(aType, dev)
			}
		})
		return true
	} else {
		ws.setState(InAlarmState, ws.armingMode, aType)
		ws.notifSystem.NotifyAlarm(aType, dev)
		return true
	}
}

func (ws *WaechterSystem) ReportBatteryLevel(level float32, dev Device) {
	if level < ws.config.BatteryThresold {
		DevLog(dev, ws.log.Info()).Float32("level", level).Msg("Battery is too low! Notify!")
		ws.notifSystem.NotifyLowBattery(dev, level)
	} else {
		DevLog(dev, ws.log.Debug()).Float32("level", level).Msg("Got battery info")
	}
}

func (ws *WaechterSystem) ReportLinkQuality(link float32, dev Device) {
	if ws.IsArmed() && link < ws.config.LinkQualityThreshold {
		DevLog(dev, ws.log.Info()).Float32("link", link).Msg("Link quality is too low! Tamper alarm!")
		ws.Alarm(TamperAlarm, dev)
	} else if link < ws.config.LinkQualityThreshold {
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
	return ws.state == ArmedState || ws.state == InAlarmState || ws.state == EntryDelayState
}
