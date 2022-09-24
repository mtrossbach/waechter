package system

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"sync"
	"time"
)

type StateUpdateFunc func(state ArmState, alarmType AlarmType)

type WaechterSystem struct {
	state               State
	notificationManager *notificationManager

	stateUpdateHandlers sync.Map
}

func NewWaechterSystem() *WaechterSystem {
	system := &WaechterSystem{
		stateUpdateHandlers: sync.Map{},
		notificationManager: newNotificationManager(cfg.GetString(cSystemName)),
	}
	system.initState()
	return system
}

func (ws *WaechterSystem) initState() {
	ws.state.loadFromDisk()

	if !isValidArmState(ws.state.ArmState) {
		ws.state.ArmState = DisarmedState
	}

	if !isValidAlarmType(ws.state.Alarm) {
		ws.state.Alarm = NoAlarm
	}

	if ws.state.Alarm == EntryDelayAlarm {
		ws.setupEntryDelayTimer(BurglarAlarm, systemDevice())
	} else if ws.state.ArmState == ExitDelayState {
		ws.setupExitDelayTimer()
	} else if ws.state.ArmState == DisarmedState {
		ws.state.Alarm = NoAlarm
	} else if ws.state.IsArmed() {
		ws.checkWrongPinCount(systemDevice())
	}

	ws.setState(ws.state.ArmState, ws.state.Alarm)
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

func (ws *WaechterSystem) ArmStay(dev Device) bool {
	return ws.arm(ArmedStayState, dev)
}

func (ws *WaechterSystem) ArmAway(dev Device) bool {
	return ws.arm(ArmedAwayState, dev)
}

func (ws *WaechterSystem) arm(targetState ArmState, dev Device) bool {
	if ws.state.IsArmedOrExitDelay() {
		ws.notifyStateHandlers()
		return false
	}

	ws.state.DelayEnd = nil
	ws.state.WrongPinCount = 0

	if targetState == ArmedStayState {
		ws.setState(ArmedStayState, NoAlarm)
	} else if targetState == ArmedAwayState {
		ws.setState(ExitDelayState, NoAlarm)
		ws.setupExitDelayTimer()
	}

	return true
}

func (ws *WaechterSystem) setupExitDelayTimer() {
	if ws.state.DelayEnd == nil {
		endTime := time.Now().Add(time.Duration(cfg.GetInt(cExitDelay)) * time.Second)
		ws.state.DelayEnd = &endTime
		ws.state.writeToDisk()
	}

	d := (*ws.state.DelayEnd).Sub(time.Now())
	if d <= 0 {
		d = 1 * time.Millisecond
	}

	time.AfterFunc(d, func() {
		ws.state.DelayEnd = nil
		if ws.state.ArmState == ExitDelayState {
			ws.setState(ArmedAwayState, NoAlarm)
		}
	})
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
		ws.state.WrongPinCount += 1
		ws.checkWrongPinCount(dev)
		return false
	}
}

func (ws *WaechterSystem) checkWrongPinCount(dev Device) {
	if ws.state.WrongPinCount > cfg.GetInt(cMaxWrongPinCount) {
		ws.Alarm(BurglarAlarm, dev)
	}
}

func (ws *WaechterSystem) ForceDisarm(dev Device) {

	if ws.state.Alarm != NoAlarm && ws.state.Alarm != EntryDelayAlarm {
		ws.notificationManager.notifyRecovery(&dev)
	}
	ws.state.WrongPinCount = 0
	ws.setState(DisarmedState, NoAlarm)
}

func (ws *WaechterSystem) Alarm(aType AlarmType, dev Device) bool {
	if !ws.state.IsArmed() && aType == BurglarAlarm {
		return false
	}
	if aType == TamperAlarm && !cfg.GetBool(cTamperAlarm) {
		return false
	}

	if ws.state.ArmState == ArmedAwayState && aType == BurglarAlarm && ws.state.Alarm == NoAlarm {
		ws.state.DelayEnd = nil
		ws.setState(ws.state.ArmState, EntryDelayAlarm)
		ws.setupEntryDelayTimer(aType, dev)
	} else if aType == BurglarAlarm && ws.state.Alarm == EntryDelayAlarm {
		// do nothing
	} else {
		ws.setState(ws.state.ArmState, aType)
		ws.notificationManager.notifyAlarm(aType, &dev)
	}
	return true
}

func (ws *WaechterSystem) setupEntryDelayTimer(aType AlarmType, dev Device) {
	if ws.state.DelayEnd == nil {
		endTime := time.Now().Add(time.Duration(cfg.GetInt(cEntryDelay)) * time.Second)
		ws.state.DelayEnd = &endTime
		ws.state.writeToDisk()
	}

	d := (*ws.state.DelayEnd).Sub(time.Now())
	if d <= 0 {
		d = 1 * time.Millisecond
	}

	time.AfterFunc(d, func() {
		ws.state.DelayEnd = nil
		if ws.state.Alarm == EntryDelayAlarm {
			DInfo(&dev).Str("alarmType", string(aType)).Msg("Alarm triggered -> entry delay")
			ws.setState(ws.state.ArmState, aType)
			ws.notificationManager.notifyAlarm(aType, &dev)
		}
	})
}

func (ws *WaechterSystem) ReportBatteryLevel(level float32, dev Device) {
	if level < cfg.GetFloat32(cBatteryThreshold) {
		DInfo(&dev).Float32("battery", level).Msg("Battery is too low! Notify!")
		ws.notificationManager.notifyLowBattery(&dev, level)
	}
}

func (ws *WaechterSystem) ReportLinkQuality(link float32, dev Device) {
	if ws.state.IsArmed() && link < cfg.GetFloat32(cLinkQualityThreshold) && cfg.GetBool(cTamperAlarm) {
		DInfo(&dev).Float32("link", link).Msg("Link quality is too low! Tamper alarm!")
		ws.Alarm(TamperAlarm, dev)
	} else if link < cfg.GetFloat32(cLinkQualityThreshold) {
		DInfo(&dev).Float32("link", link).Msg("Link quality is too low! Notify!")
		ws.notificationManager.notifyLowLinkQuality(&dev, link)
	}
}

func (ws *WaechterSystem) GetArmState() ArmState {
	return ws.state.ArmState
}

func (ws *WaechterSystem) GetAlarmType() AlarmType {
	return ws.state.Alarm
}

func (ws *WaechterSystem) setState(state ArmState, alarmType AlarmType) {
	ws.state.ArmState = state
	ws.state.Alarm = alarmType
	ws.state.writeToDisk()
	log.Info().Str("armState", string(state)).Str("alarm", string(alarmType)).Msg("System state updated.")
	ws.notifyStateHandlers()
}

func (ws *WaechterSystem) notifyStateHandlers() {
	ws.stateUpdateHandlers.Range(func(_, value any) bool {
		handler := value.(StateUpdateFunc)
		handler(ws.state.ArmState, ws.state.Alarm)
		return true
	})
}
