package zigbee2mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/internal/wslice"
	"github.com/mtrossbach/waechter/system"
	"github.com/mtrossbach/waechter/system/arm"
	"github.com/mtrossbach/waechter/system/device"
	"sync"
	"time"
)

type Connector struct {
	conf             config.Zigbee2MqttConfig
	ctrl             system.Controller
	conn             *connection
	availableDevices sync.Map //map[device.Id]device.Spec
	activeDevices    sync.Map //map[device.Id]nil
	connected        bool
}

func NewConnector(configuration config.Zigbee2MqttConfig) (*Connector, error) {
	if len(configuration.Id) == 0 {
		return nil, errors.New("no id")
	}

	if len(configuration.Url) == 0 {
		return nil, errors.New("no url")
	}

	if len(configuration.BaseTopic) == 0 {
		configuration.BaseTopic = "zigbee2mqtt"
	}

	return &Connector{
		conf:             configuration,
		conn:             newConnection(configuration),
		availableDevices: sync.Map{},
		activeDevices:    sync.Map{},
		connected:        false,
	}, nil
}

func (c *Connector) Setup(controller system.Controller) {
	c.ctrl = controller
	c.conn.OnConnect = func(conn *connection) {
		log.Info().Str("id", c.conf.Id).Str("url", c.conf.Url).Msg("Connected to Zigbee2Mqtt broker")
		c.connected = true
		c.ctrl.OperationalStateChanged(c)
	}

	c.conn.OnConnectionLost = func(conn *connection, err error) {
		log.Error().Err(err).Str("id", c.conf.Id).Str("url", c.conf.Url).Msg("Connection to Zigbee2Mqtt broker lost. Reconnecting ...")
		c.connected = false
		c.ctrl.OperationalStateChanged(c)
	}

	log.Debug().Str("id", c.conf.Id).Str("url", c.conf.Url).Msg("Connecting to Zigbee2Mqtt broker...")
	c.conn.Subscribe("bridge/devices", c.handleNewDeviceList)
	c.conn.Subscribe("bridge/event", c.handleDeviceEvent)
	c.conn.Connect()

}

func (c *Connector) Teardown() {
	c.activeDevices.Range(func(key, _ any) bool {
		c.DeactivateDevice(key.(device.Id))
		return true
	})
	c.conn.Disconnect()
	c.ctrl = nil
}

func (c *Connector) Id() string {
	return c.conf.Id
}

func (c *Connector) DisplayName() string {
	return "Zigbee2Mqtt Connector"
}

func (c *Connector) Operational() bool {
	return c.connected
}

func (c *Connector) EnumerateDevices() []device.Spec {
	var data []device.Spec

	c.availableDevices.Range(func(_, value any) bool {
		data = append(data, value.(device.Spec))
		return true
	})

	return data
}

func (c *Connector) ActivateDevice(id device.Id) error {
	_, found := c.activeDevices.Load(id)
	if found {
		return nil
	}

	_, found = c.availableDevices.Load(id)
	if !found {
		return errors.New("device not found")
	}

	c.activeDevices.Store(id, nil)
	c.conn.Subscribe(id.Entity(), c.deviceMessageHandler(id))
	c.ctrl.DeviceAvailable(id)
	return nil
}

func (c *Connector) DeactivateDevice(id device.Id) error {
	_, found := c.activeDevices.Load(id)
	if found {
		return errors.New("device not found")
	}

	c.activeDevices.Delete(id)
	c.conn.Unsubscribe(id.Entity())
	c.ctrl.DeviceUnavailable(id)
	return nil
}

func (c *Connector) ControlActor(id device.Id, actor device.Actor, value any) bool {
	switch actor {
	case device.StateActor:
		c.sendPayload(id, newArmModePayload(c.ctrl.SystemState(), nil))
		return true

	case device.AlarmActor:
		c.sendPayload(id, newWarningPayload(c.ctrl.SystemState().Alarm))
		time.AfterFunc(100*time.Millisecond, func() {
			// Resend after 100ms because some sirens do not correctly process payloads during active alarms
			c.sendPayload(id, newWarningPayload(c.ctrl.SystemState().Alarm))
		})
		return true

	default:
		log.Error().Str("device", string(id)).Str("actor", string(actor)).Interface("value", value).Msg("Unknown actor type")
	}

	return false
}

func (c *Connector) sendPayload(id device.Id, payload any) {
	c.conn.Publish(fmt.Sprintf("%v/set", id.Entity()), payload)
}

func (c *Connector) deviceMessageHandler(id device.Id) MessageHandler {
	return func(msg mqtt.Message) {
		var data map[string]any
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Error().Str("device", string(id)).Str("payload", string(msg.Payload())).Msg("Could not parse device data")
			return
		}

		spec, ok := c.availableDevices.Load(id)
		if !ok {
			log.Error().Str("device", string(id)).Str("payload", string(msg.Payload())).Msg("Could not process device data for unknown device")
		}

		actionDone := false
		for _, s := range spec.(device.Spec).Sensors {
			switch s {
			case device.MotionSensor:
				if v := extract[bool](data, "occupancy"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.MotionSensorValue{Motion: *v})
				}
			case device.ContactSensor:
				if v := extract[bool](data, "contact"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.ContactSensorValue{Contact: *v})
				}
			case device.SmokeSensor:
				if v := extract[bool](data, "smoke"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.SmokeSensorValue{Smoke: *v})
				}
			case device.BatteryWarningSensor:
				if v := extract[bool](data, "battery_low"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.BatteryWarningSensorValue{BatteryWarning: *v})
				}
			case device.TamperSensor:
				if v := extract[bool](data, "tamper"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.TamperSensorValues{Tamper: *v})
				}
			case device.BatteryLevelSensor:
				if v := extract[int](data, "battery"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.BatteryLevelSensorValue{BatteryLevel: float32(*v) / float32(100)})
				}
			case device.LinkQualitySensor:
				if v := extract[int](data, "linkquality"); v != nil {
					c.ctrl.DeliverSensorValue(id, s, device.LinkQualitySensorValue{LinkQuality: float32(*v) / float32(255)})
				}
			case device.ArmingSensor, device.DisarmingSensor, device.PanicSensor:
				if actionDone {
					continue
				} else {
					actionDone = true
				}
				if v := extract[string](data, "action"); v != nil && len(*v) > 0 {
					transactionId := extract[float64](data, "action_transaction")
					c.sendPayload(id, armModePayload{ArmMode: armMode{
						Mode:        *v,
						Transaction: transactionId,
					}})
					c.ctrl.DeliverSensorValue(id, s, device.PanicSensorValue{Panic: *v == "panic"})

					if *v == "arm_all_zones" {
						c.ctrl.DeliverSensorValue(id, device.ArmingSensor, device.ArmingSensorValue{ArmMode: arm.All})
					} else if *v == "arm_day_zones" || *v == "arm_night_zones" {
						c.ctrl.DeliverSensorValue(id, device.ArmingSensor, device.ArmingSensorValue{ArmMode: arm.Perimeter})
					} else if *v == "disarm" {
						if code := extract[string](data, "action_code"); code != nil {
							c.ctrl.DeliverSensorValue(id, device.ArmingSensor, device.DisarmingSensorValue{Pin: *code})
						}
					}
				}
			}
		}

	}
}

func (c *Connector) handleDeviceEvent(msg mqtt.Message) {
	var deviceEvent DeviceEvent
	if err := json.Unmarshal(msg.Payload(), &deviceEvent); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse Zigbee device event!")
		return
	}

	if deviceEvent.Type == "device_announce" && len(deviceEvent.Data.FriendlyName) > 0 {
		id := device.NewId(c.Id(), deviceEvent.Data.FriendlyName)
		c.ctrl.DeviceAvailable(id)
	}
}

func (c *Connector) handleNewDeviceList(msg mqtt.Message) {

	var newDevices []Z2MDeviceInfo
	if err := json.Unmarshal(msg.Payload(), &newDevices); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse Zigbee device payload!")
		return
	}

	relevantDevices := make(map[string]Z2MDeviceInfo)
	for _, d := range newDevices {
		if (d.Type == "EndDevice" || d.Type == "Router") && d.Supported {
			relevantDevices[d.FriendlyName] = d
		}
	}

	c.availableDevices = sync.Map{}

	for _, d := range relevantDevices {
		spec := c.specFromDeviceInfo(d)
		if spec.IsRelevant() {
			c.availableDevices.Store(spec.Id, spec)
		}
	}

	c.activeDevices.Range(func(key, _ any) bool {
		_, found := c.availableDevices.Load(key)
		if !found {
			c.DeactivateDevice(key.(device.Id))
		}
		return true
	})

	c.ctrl.DeviceListUpdated(c)
}

func (c *Connector) specFromDeviceInfo(info Z2MDeviceInfo) device.Spec {
	spec := device.Spec{
		Id:          device.NewId(c.Id(), info.FriendlyName),
		DisplayName: info.FriendlyName,
		Sensors:     []device.Sensor{},
		Actors:      []device.Actor{},
	}

	var exposes []string
	for _, e := range info.Definition.Exposes {
		exposes = append(exposes, e.Property)
	}

	if wslice.ContainsAll(exposes, []string{"action_code", "action"}) {
		spec.Sensors = append(spec.Sensors, device.ArmingSensor, device.DisarmingSensor, device.PanicSensor)
		spec.Actors = append(spec.Actors, device.StateActor)
	}

	if wslice.ContainsAll(exposes, []string{"warning"}) {
		spec.Actors = append(spec.Actors, device.AlarmActor)
	}

	if wslice.ContainsAll(exposes, []string{"contact"}) {
		spec.Sensors = append(spec.Sensors, device.ContactSensor)
	}

	if wslice.ContainsAll(exposes, []string{"smoke"}) {
		spec.Sensors = append(spec.Sensors, device.SmokeSensor)
	}

	if wslice.ContainsAll(exposes, []string{"occupancy"}) {
		spec.Sensors = append(spec.Sensors, device.MotionSensor)
	}

	if wslice.ContainsAll(exposes, []string{"battery"}) {
		spec.Sensors = append(spec.Sensors, device.BatteryLevelSensor)
	} else if wslice.ContainsAll(exposes, []string{"battery_low"}) {
		spec.Sensors = append(spec.Sensors, device.BatteryWarningSensor)
	}

	if wslice.ContainsAll(exposes, []string{"tamper"}) {
		spec.Sensors = append(spec.Sensors, device.TamperSensor)
	}

	if wslice.ContainsAll(exposes, []string{"linkquality"}) {
		spec.Sensors = append(spec.Sensors, device.LinkQualitySensor)
	}

	return spec
}
