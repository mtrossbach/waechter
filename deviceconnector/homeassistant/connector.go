package homeassistant

import (
	"errors"
	"github.com/mtrossbach/waechter/deviceconnector/homeassistant/connection"
	"github.com/mtrossbach/waechter/deviceconnector/homeassistant/msgs"
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
	"github.com/mtrossbach/waechter/system/device"
	"strconv"
	"strings"
	"sync"
)

type Connector struct {
	conf             config.HomeAssistantConfig
	ctrl             system.Controller
	conn             *connection.Connection
	availableDevices sync.Map //map[device.Id]assembledDevice
	activeDevices    sync.Map //map[device.Id]nil
	connected        bool
}

func NewConnector(configuration config.HomeAssistantConfig) (*Connector, error) {
	if len(configuration.Id) == 0 {
		return nil, errors.New("no id")
	}

	if len(configuration.Url) == 0 {
		return nil, errors.New("no url")
	}

	if len(configuration.Token) == 0 {
		return nil, errors.New("no token")
	}

	return &Connector{
		conf:             configuration,
		conn:             connection.NewConnection(configuration.Url, configuration.Token),
		availableDevices: sync.Map{},
		activeDevices:    sync.Map{},
		connected:        false,
	}, nil
}

func (c *Connector) Setup(controller system.Controller) {
	c.ctrl = controller

	log.Debug().Str("url", c.conf.Url).Msg("Connecting to HomeAssistant...")

	c.conn.OnConnect = func(conn *connection.Connection) {
		log.Info().Str("id", c.conf.Id).Str("url", c.conf.Url).Msg("Connected to HomeAssistant")
		c.connected = true
		c.ctrl.OperationalStateChanged(c)
		c.updateDeviceList()
	}

	c.conn.OnConnectionLost = func(conn *connection.Connection, err error) {
		log.Error().Err(err).Str("id", c.conf.Id).Str("url", c.conf.Url).Msg("Connection to HomeAssistant lost. Reconnecting ...")
		c.connected = false
		c.ctrl.OperationalStateChanged(c)
	}
	c.conn.Connect()
}

func (c *Connector) updateDeviceList() {
	st, err := c.getStates()
	if err != nil {
		log.Error().Err(err).Msg("Could not request states from HomeAssistant")
		return
	}

	devs := make(map[string]assembledDevice)

	for _, s := range st.Result {
		prefix := entityPrefix(s.EntityID)
		dev, ok := devs[prefix]
		if !ok {
			dev = assembledDevice{}
		}

		switch s.Attributes.DeviceClass {
		case "motion":
			dev.sensors[device.MotionSensor] = s.EntityID
		case "opening", "door", "window", "garage_door":
			dev.sensors[device.ContactSensor] = s.EntityID
		case "smoke":
			dev.sensors[device.SmokeSensor] = s.EntityID
		case "battery":
			if strings.HasPrefix(s.EntityID, "binary_sensor") {
				dev.sensors[device.BatteryWarningSensor] = s.EntityID
			} else {
				dev.sensors[device.BatteryLevelSensor] = s.EntityID
			}
		case "tamper":
			dev.sensors[device.TamperSensor] = s.EntityID
		}

		dev.spec = dev.generateSpec(c.Id())
		devs[prefix] = dev
	}

	c.activeDevices.Range(func(key, _ any) bool {
		_, found := c.availableDevices.Load(key)
		if !found {
			_ = c.DeactivateDevice(key.(device.Id))
		}
		return true
	})

	c.ctrl.DeviceListUpdated(c)
}

func (c *Connector) getStates() (msgs.StateResult, error) {
	var result msgs.StateResult
	err := c.conn.Command(&msgs.BaseMessage{Type: msgs.GetStates}, &result)
	return result, err
}

func (c *Connector) Teardown() {

}

func (c *Connector) Id() string {
	return c.conf.Id
}

func (c *Connector) DisplayName() string {
	return "HomeAssistant Connector"
}

func (c *Connector) Operational() bool {
	return c.connected
}

func (c *Connector) EnumerateDevices() []device.Spec {
	var data []device.Spec

	c.availableDevices.Range(func(_, value any) bool {
		data = append(data, value.(assembledDevice).spec)
		return true
	})

	return data
}

func (c *Connector) ActivateDevice(id device.Id) error {
	_, found := c.activeDevices.Load(id)
	if found {
		return nil
	}

	dev, found := c.availableDevices.Load(id)
	if !found {
		return errors.New("device not found")
	}

	if d, ok := dev.(assembledDevice); ok {
		for sensor, entityId := range d.sensors {
			if err := c.conn.SubscribeStateEvents(entityId, c.deviceEventHandler(id, sensor)); err != nil {
				log.Error().Err(err).Str("device", string(id)).Str("entityId", entityId).Msg("Could not setup state trigger for device")
			}
		}
	}

	return nil
}

func (c *Connector) DeactivateDevice(id device.Id) error {
	_, found := c.activeDevices.Load(id)
	if found {
		return errors.New("device not found")
	}

	dev, found := c.availableDevices.Load(id)
	if !found {
		return errors.New("device not found")
	}

	if d, ok := dev.(assembledDevice); ok {
		for _, entityId := range d.sensors {
			if err := c.conn.UnsubscribeStateEvents(entityId); err != nil {
				log.Error().Err(err).Str("device", string(id)).Str("entityId", entityId).Msg("Could not uninstall state trigger for device")
			}
		}
	}

	return nil
}

func (c *Connector) ControlActor(id device.Id, actor device.Actor, value any) bool {
	//Not yet supported
	return true
}

func (c *Connector) deviceEventHandler(id device.Id, sensor device.Sensor) connection.StateEventHandler {
	return func(entityId string, event msgs.EventResponse) {
		switch sensor {
		case device.MotionSensor:
			c.ctrl.DeliverSensorValue(id, sensor, device.MotionSensorValue{Motion: event.Event.Variables.Trigger.ToState.State == "on"})

		case device.ContactSensor:
			c.ctrl.DeliverSensorValue(id, sensor, device.ContactSensorValue{Contact: event.Event.Variables.Trigger.ToState.State == "on"})

		case device.SmokeSensor:
			c.ctrl.DeliverSensorValue(id, sensor, device.SmokeSensorValue{Smoke: event.Event.Variables.Trigger.ToState.State == "on"})

		case device.BatteryLevelSensor:
			level, err := strconv.Atoi(event.Event.Variables.Trigger.ToState.State)
			if err != nil {
				log.Error().Err(err).Str("state", event.Event.Variables.Trigger.ToState.State).Msg("Could not parse battery level")
			} else {
				c.ctrl.DeliverSensorValue(id, sensor, device.BatteryLevelSensorValue{BatteryLevel: float32(level) / 100.0})
			}

		case device.BatteryWarningSensor:
			c.ctrl.DeliverSensorValue(id, sensor, device.BatteryWarningSensorValue{BatteryWarning: event.Event.Variables.Trigger.ToState.State == "on"})
		case device.TamperSensor:
			c.ctrl.DeliverSensorValue(id, sensor, device.TamperSensorValues{Tamper: event.Event.Variables.Trigger.ToState.State == "on"})

		}
	}
}
