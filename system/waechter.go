package system

import (
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/internal/wslice"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/mtrossbach/waechter/system/arm"
	"github.com/mtrossbach/waechter/system/device"
	"github.com/mtrossbach/waechter/system/zone"
	"sync"
	"time"
)

type Waechter struct {
	name             string
	state            State
	zones            map[zone.Id]*zone.Zone
	devices          map[device.Id]*device.Device
	deviceConnectors []DeviceConnector
	wrongPinCount    int
	noteMgr          *notificationManager

	entryTimers          sync.Map
	unavailabilityTimers sync.Map
}

func NewWaechter() *Waechter {
	w := Waechter{
		state:                State{},
		zones:                nil,
		devices:              nil,
		wrongPinCount:        0,
		deviceConnectors:     []DeviceConnector{},
		noteMgr:              newNotificationManager(),
		entryTimers:          sync.Map{},
		unavailabilityTimers: sync.Map{},
	}

	w.loadState()
	w.loadZones()
	w.loadDevices()
	return &w
}

func (w *Waechter) AddDeviceConnector(connector DeviceConnector) {
	w.deviceConnectors = append(w.deviceConnectors, connector)
	connector.Setup(w)
}

func (w *Waechter) RemoveDeviceConnector(id string) {
	connector, i := wslice.FilterOne[DeviceConnector](w.deviceConnectors,
		func(i DeviceConnector) bool { return i.Id() == id })
	if connector != nil {
		(*connector).Teardown()
		w.deviceConnectors = wslice.Remove[DeviceConnector](w.deviceConnectors, i)
	}
}

func (w *Waechter) AddNotificationAdapter(adapter NotificationAdapter) {
	w.noteMgr.AddAdapter(adapter)
}

func (w *Waechter) loadState() {
	w.setArmMode(arm.Disarmed)
	w.setAlarm(alarm.None)
}

func (w *Waechter) loadZones() {
	w.zones = make(map[zone.Id]*zone.Zone)
	for _, zc := range config.ZoneConfigs() {
		z := zone.ZoneFromConfig(zc)
		w.zones[z.Id] = &z
	}
}

func (w *Waechter) loadDevices() {
	w.devices = make(map[device.Id]*device.Device)
	for _, dc := range config.DeviceConfigs() {
		d := device.DeviceFromConfig(dc)
		w.devices[d.Id] = &d
	}
	w.devices[systemDeviceId] = systemDevice()
}

func (w *Waechter) zoneForDeviceId(id device.Id) zone.Zone {
	z, ok := w.zones[w.devices[id].Zone]
	if !ok {
		return zone.SubstitutionZone(w.name, w.state.Armed())
	}
	return *z
}

func (w *Waechter) deviceConnectorForId(id string) DeviceConnector {
	c, _ := wslice.FilterOne[DeviceConnector](w.deviceConnectors, func(i DeviceConnector) bool { return i.Id() == id })
	return *c
}

func (w *Waechter) DeliverSensorValue(id device.Id, sensor device.Sensor, value any) bool {
	oldValue := w.devices[id].State[sensor]
	w.devices[id].State[sensor] = value

	if oldValue != nil && oldValue == value {
		return false
	}

	z := w.zoneForDeviceId(id)

	if v, ok := value.(device.MotionSensorValue); ok {
		if z.Armed && v.Motion {
			if !(z.Delayed && w.isDuringExitDelay()) {
				w.alarm(id, alarm.Burglar, z.Delayed)
			}
		}

	} else if v, ok := value.(device.ContactSensorValue); ok {
		if z.Armed && !v.Contact {
			if !(z.Delayed && w.isDuringExitDelay()) {
				w.alarm(id, alarm.Burglar, z.Delayed)
			}
		}

	} else if v, ok := value.(device.SmokeSensorValue); ok {
		if v.Smoke {
			w.alarm(id, alarm.Fire, false)
		}

	} else if v, ok := value.(device.PanicSensorValue); ok {
		if v.Panic {
			w.alarm(id, alarm.Panic, false)
		}

	} else if v, ok := value.(device.BatteryWarningSensorValue); ok {
		if v.BatteryWarning {
			w.noteMgr.NotifyLowBattery(w.specForDeviceId(id), w.zoneForDeviceId(id), 0)
		}

	} else if v, ok := value.(device.TamperSensorValues); ok {
		if v.Tamper {
			if (z.Armed && config.TamperAlarmWhileArmed()) || (!z.Armed && config.TamperAlarmWhileDisarmed()) {
				w.alarm(id, alarm.Tamper, false)
			}
		}

	} else if v, ok := value.(device.BatteryLevelSensorValue); ok {
		if v.BatteryLevel < config.BatteryLevelThreshold() {
			w.noteMgr.NotifyLowBattery(w.specForDeviceId(id), w.zoneForDeviceId(id), v.BatteryLevel)
		}

	} else if v, ok := value.(device.LinkQualitySensorValue); ok {
		if v.LinkQuality < config.LinkQualityThreshold() {
			w.noteMgr.NotifyLowLinkQuality(w.specForDeviceId(id), w.zoneForDeviceId(id), v.LinkQuality)
		}

	} else if v, ok := value.(device.ArmingSensorValue); ok {
		if v.ArmMode == arm.Disarmed {
			return false
		}
		return w.arm(id, v.ArmMode)

	} else if v, ok := value.(device.DisarmingSensorValue); ok {
		return w.disarm(id, v.Pin)

	} else {
		log.Error().Str("device", string(id)).Interface("value", value).Msg("Unknown sensor value received")
		return false
	}

	return true
}

func (w *Waechter) isDuringExitDelay() bool {
	exitDelay := time.Duration(config.ExitDelay()) * time.Second
	return w.state.Armed() && time.Now().Sub(w.state.armModeUpdated) < exitDelay
}

func (w *Waechter) alarm(id device.Id, alarmType alarm.Type, delayedZone bool) {
	if alarmType == alarm.Burglar && delayedZone && (w.state.Alarm == alarm.None || w.state.Alarm == alarm.EntryDelay) {
		w._alarm(id, alarm.EntryDelay)
		t, ok := w.entryTimers.Load(id)
		if !ok {
			t = time.AfterFunc(time.Duration(config.EntryDelay())*time.Second, func() {
				w.entryTimers.Delete(id)
				if w.zoneForDeviceId(id).Armed {
					w._alarm(id, alarmType)
				}
			})
			w.entryTimers.Store(id, t)
		}

	} else {
		w._alarm(id, alarmType)
	}
}

func (w *Waechter) specForDeviceId(id device.Id) device.Spec {
	d, ok := w.devices[id]
	if !ok {
		d = systemDevice()
	}
	return (*d).Spec
}

func (w *Waechter) _alarm(id device.Id, alarm alarm.Type) {
	w.setAlarm(alarm)
	w.noteMgr.NotifyAlarm(alarm, w.specForDeviceId(id), w.zoneForDeviceId(id))
}

func (w *Waechter) arm(id device.Id, mode arm.Mode) bool {
	if w.state.Armed() || mode == arm.Disarmed {
		return false
	}
	if mode == arm.Disarmed {
		mode = arm.All
	}
	w.wrongPinCount = 0
	w.setArmMode(mode)
	return true
}

func (w *Waechter) disarm(id device.Id, enteredPin string) bool {
	persons := config.Persons()
	person, _ := wslice.FilterOne(persons, func(p config.Person) bool { return p.Pin == enteredPin })

	if person != nil {
		w.wrongPinCount = 0
		if w.state.Alarm != alarm.None {
			w.noteMgr.NotifyRecovery(w.specForDeviceId(id), w.zoneForDeviceId(id))
		}
		log.Info().Str("name", person.Name).Msg("Disarmed by pin")
		w.setAlarm(alarm.None)
		w.setArmMode(arm.Disarmed)
		w.entryTimers.Range(func(key, value any) bool {
			t := value.(*time.Timer)
			t.Stop()
			w.entryTimers.Delete(key)
			return true
		})
		return true
	} else {
		w.wrongPinCount += 1
		log.Info().Str("device", string(id))
		if w.wrongPinCount > config.MaxWrongPinCount() {
			w.alarm(id, alarm.TamperPin, false)
		}
		return false
	}
}

func (w *Waechter) SystemState() State {
	return w.state
}

func (w *Waechter) DeviceListUpdated(system DeviceConnector) {
	if system == nil {
		return
	}
	deviceSpecs := system.EnumerateDevices()
	log.Info().Str("connector", system.DisplayName()).Str("id", system.Id()).Msg("Received new device list:")
	for _, s := range deviceSpecs {
		if ad, ok := w.devices[s.Id]; ok {
			ad.Spec = s
		}
		var sensors []string
		var actors []string
		for _, ss := range s.Sensors {
			sensors = append(sensors, string(ss))
		}
		for _, sa := range s.Actors {
			actors = append(actors, string(sa))
		}
		log.Info().Str("id", string(s.Id)).Str("displayName", s.DisplayName).Str("vendor", s.Vendor).Str("model", s.Model).Strs("sensors", sensors).Strs("actors", actors).Msg("\t- Device detected")
	}

	for _, d := range w.devices {
		if !d.Active && d.Id.Prefix() == system.Id() {
			err := system.ActivateDevice(d.Id)
			if err != nil {
				device.DError(d).Err(err).Msg("✗ Could not activate device")
			} else {
				device.DInfo(d).Msg("✓ Device active")
			}
		}
	}
}

func (w *Waechter) OperationalStateChanged(connector DeviceConnector) {
	if !connector.Operational() && config.DeviceSystemFaultAlarm() && w.state.Armed() {
		time.AfterFunc(time.Duration(config.DeviceSystemFaultDelay())*time.Second, func() {
			if !connector.Operational() && w.state.Armed() {
				w.alarm(systemDeviceId, alarm.Tamper, false)
			}
		})
	}
}

func (w *Waechter) DeviceUnavailable(id device.Id) {
	d, ok := w.devices[id]
	if ok {
		d.Active = false
	}

	z := w.zoneForDeviceId(id)
	if z.Armed {
		w.alarm(id, alarm.Tamper, false)
	}
}

func (w *Waechter) DeviceAvailable(id device.Id) {
	d, ok := w.devices[id]
	if ok {
		d.Active = true
		w.updateActor(id, device.StateActor, w.state.stateActorPayload())
		w.updateActor(id, device.AlarmActor, w.state.alarmActorPayload())
	}
}

func (w *Waechter) setArmMode(mode arm.Mode) {
	if w.state.ArmMode != mode {
		w.state.ArmMode = mode
		w.state.armModeUpdated = time.Now()

		w.syncZones()

		w.updateActors(device.StateActor, w.state.stateActorPayload())

		l := log.Info().Str("mode", string(mode))
		if w.state.Armed() {
			l = l.Int("exitDelay", config.ExitDelay())
		}
		l.Msg("System mode changed")

	}
}

func (w *Waechter) setAlarm(a alarm.Type) {
	if w.state.Alarm != a {
		w.state.Alarm = a

		w.updateActors(device.StateActor, w.state.stateActorPayload())
		w.updateActors(device.AlarmActor, w.state.alarmActorPayload())

		l := log.Info().Str("alarm", string(a))
		if a == alarm.EntryDelay {
			l = l.Int("entryDelay", config.EntryDelay())
		}
		l.Msg("Alarm changed")
	}
}

func (w *Waechter) updateActor(id device.Id, actor device.Actor, payload any) {
	if d, ok := w.devices[id]; ok && d != nil && wslice.Contains(d.Spec.Actors, actor) {
		if c := w.deviceConnectorForId(d.Id.Prefix()); c != nil {
			c.ControlActor(d.Id, actor, payload)
		}
	}
}

func (w *Waechter) updateActors(actor device.Actor, payload any) {
	for i := range w.devices {
		w.updateActor(i, actor, payload)
	}
}

func (w *Waechter) syncZones() {
	for _, z := range w.zones {
		if z.Perimeter {
			z.Armed = w.state.Armed()
		} else {
			if w.state.Armed() && w.state.ArmMode != arm.Perimeter {
				z.Armed = true
			} else {
				z.Armed = false
			}
		}
	}
}
