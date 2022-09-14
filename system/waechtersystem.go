package system

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"sync"
	"time"
)

type StateUpdateFunc func(state State, armingMode ArmingMode, alarmType AlarmType)

type WaechterSystem struct {
	state               State
	armingMode          ArmingMode
	alarmType           AlarmType
	wrongPinCount       int
	notificationManager *notificationManager

	stateUpdateHandlers sync.Map
}

func NewWaechterSystem() *WaechterSystem {
	system := &WaechterSystem{
		wrongPinCount:       0,
		stateUpdateHandlers: sync.Map{},
		notificationManager: newNotificationManager(cfg.GetString(cSystemName)),
	}
	system.initState()
	return system
}

func (ws *WaechterSystem) initState() {
	cfgState, cfgArmingMode, cfgAlarmType := loadState()

	if isValidState(cfgState) && isValidArmingMode(cfgArmingMode) && isValidAlarmType(cfgAlarmType) {
		if cfgState == EntryDelayState {
			cfgState = ArmedState
			cfgAlarmType = BurglarAlarm
		} else if cfgState == ArmingState {
			cfgState = ArmedState
		}
		ws.setState(cfgState, cfgArmingMode, cfgAlarmType)
	} else {
		ws.setState(DisarmedState, AwayMode, NoAlarm)
	}
}

func (ws *WaechterSystem) SubscribeStateUpdate(id interface{}, fun StateUpdateFunc) {
	ws.stateUpdateHandlers.Store(id, fun)
}

func (ws *WaechterSystem) UnsubscribeStateUpdate(id interface{}) {
	ws.stateUpdateHandlers.Delete(id)
}

func (ws *WaechterSystem) AddNotificationAdapter(adapter NotificationAdapter) {
	ws.notificationManager.AddAdapter(adapter)
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
		ws.notifyStateHandlers()
		return false
	}
}

func (ws *WaechterSystem) Disarm(enteredPin string, dev Device) bool {
	pinOk := false
	disarmPins := cfg.GetStrings(cDisarmPins)
	for _, pin := range disarmPins {
		if pin == enteredPin {
			pinOk = true
			break
		}
	}

	if pinOk {
		ws.ForceDisarm(dev)
		return true
	} else {
		ws.wrongPinCount += 1
		if ws.wrongPinCount > cfg.GetInt(cMaxWrongPinCount) {
			ws.Alarm(BurglarAlarm, dev)
		}
		return false
	}
}

func (ws *WaechterSystem) ForceDisarm(dev Device) {
	ws.wrongPinCount = 0
	if ws.alarmType != NoAlarm {
		ws.notificationManager.notifyRecovery(&dev)
	}
	ws.setState(DisarmedState, ws.armingMode, NoAlarm)
}

func (ws *WaechterSystem) Alarm(aType AlarmType, dev Device) bool {
	if (ws.state != ArmedState) && aType == BurglarAlarm {
		return false
	}
	if aType == TamperAlarm && !cfg.GetBool(cTamperAlarm) {
		return false
	}

	if ws.state == ArmedState && aType == BurglarAlarm && ws.alarmType == NoAlarm {
		ws.setState(EntryDelayState, ws.armingMode, ws.alarmType)
		time.AfterFunc(time.Duration(cfg.GetInt(cEntryDelay))*time.Second, func() {
			if ws.state == EntryDelayState {
				DInfo(&dev).Str("alarmType", string(aType)).Msg("Alarm triggered -> entry delay")
				ws.setState(ArmedState, ws.armingMode, aType)
				ws.notificationManager.notifyAlarm(aType, &dev)
			}
		})
		return true
	} else {
		ws.setState(ws.state, ws.armingMode, aType)
		ws.notificationManager.notifyAlarm(aType, &dev)
		return true
	}
}

func (ws *WaechterSystem) ReportBatteryLevel(level float32, dev Device) {
	if level < cfg.GetFloat32(cBatteryThreshold) {
		DInfo(&dev).Float32("battery", level).Msg("Battery is too low! Notify!")
		ws.notificationManager.notifyLowBattery(&dev, level)
	}
}

func (ws *WaechterSystem) ReportLinkQuality(link float32, dev Device) {
	if ws.IsArmed() && link < cfg.GetFloat32(cLinkQualityThreshold) && cfg.GetBool(cTamperAlarm) {
		DInfo(&dev).Float32("link", link).Msg("Link quality is too low! Tamper alarm!")
		ws.Alarm(TamperAlarm, dev)
	} else if link < cfg.GetFloat32(cLinkQualityThreshold) {
		DInfo(&dev).Float32("link", link).Msg("Link quality is too low! Notify!")
		ws.notificationManager.notifyLowLinkQuality(&dev, link)
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
	log.Info().Str("state", string(state)).Str("armingMode", string(mode)).Str("alarmType", string(alarmType)).Msg("System state updated")
	ws.notifyStateHandlers()
	saveState(state, mode, alarmType)
}

func (ws *WaechterSystem) notifyStateHandlers() {
	ws.stateUpdateHandlers.Range(func(_, value any) bool {
		handler := value.(StateUpdateFunc)
		handler(ws.state, ws.armingMode, ws.alarmType)
		return true
	})
}

func (ws *WaechterSystem) IsArmed() bool {
	return ws.state == ArmedState || ws.state == EntryDelayState
}
