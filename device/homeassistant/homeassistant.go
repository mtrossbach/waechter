package homeassistant

import (
	dd "github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/homeassistant/connector"
	"github.com/mtrossbach/waechter/device/homeassistant/driver"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/internal/wstring"
	"github.com/mtrossbach/waechter/system"
	"strings"
	"sync"
	"time"
)

const namespace string = "ha"

type homeassistant struct {
	connector  *connector.Connector
	devices    sync.Map
	connection uint64
}

func New() *homeassistant {
	return &homeassistant{
		connection: 0,
	}
}

func (ha *homeassistant) Start(controller dd.SystemController) {
	ha.devices = sync.Map{}
	ha.connector = connector.NewConnector()
	log.Debug().Str("uri", cfg.GetString(cURL)).Msg("Connecting to HomeAssistant...")
	err := ha.connector.Connect(cfg.GetString(cURL), cfg.GetString(cToken), ha.disconnectedHandler(controller))
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to HomeAssistant. Retrying in a few seconds...")
		ha.reconnect(controller)
		return
	}
	log.Info().Str("uri", cfg.GetString(cURL)).Msg("Connected to HomeAssistant.")
	go ha.testConnection(ha.connection)
	st, err := ha.getStates()
	if err != nil {
		log.Error().Err(err).Msg("Could not request states from HomeAssistant")
		return
	}

	deviceConfig := ha.getDeviceConfig()

	for _, s := range st.Result {
		config, ok := deviceConfig[s.EntityID]
		if ok {
			delete(deviceConfig, s.EntityID)
		} else if !cfg.GetBool(cAutoDeviceDiscovery) {
			continue
		}

		dev := ha.deviceFromSensorStateAndConfig(&s, &config)
		if dev != nil {
			var batteryId string
			var tamperId string
			if ok && len(config.Battery) > 0 {
				batteryId = config.Battery
			} else {
				prefix := strings.Replace(entityPrefix(dev.Id), "binary_sensor.", "sensor.", 1)
				batteryId = st.GetEntityIdWithPrefixAndType(prefix, "battery")
			}
			if ok && len(config.Tamper) > 0 {
				tamperId = config.Tamper
			} else {
				prefix := entityPrefix(dev.Id)
				tamperId = st.GetEntityIdWithPrefixAndType(prefix, "tamper")
			}

			ha.setupDevice(*dev, batteryId, tamperId, controller)

		}
	}

	for k := range deviceConfig {
		log.Error().Str("_id", k).Msg("Cannot add device because it is not available in HomeAssistant!")
	}

}

func (ha *homeassistant) deviceFromSensorStateAndConfig(state *msgs.SensorState, config *devicesConfig) *system.Device {
	if config == nil && state != nil {
		return ha.deviceFromSensorState(*state)
	}

	if len(config.Type) > 0 && len(config.Id) > 0 {
		return &system.Device{
			Namespace: namespace,
			Id:        config.Id,
			Name:      wstring.StrDef(config.Name, config.Id),
			Type:      config.Type,
		}
	} else if state != nil {
		dev := ha.deviceFromSensorState(*state)
		if dev != nil {
			dev.Name = wstring.StrDef(config.Name, dev.Name)
			dev.Type = system.DeviceType(wstring.StrDef(string(config.Type), string(dev.Type)))
		}
		return dev
	} else {
		return nil
	}
}

func (ha *homeassistant) deviceFromSensorState(s msgs.SensorState) *system.Device {
	if s.Attributes.DeviceClass == "motion" && (s.Attributes.MotionValid == nil || *s.Attributes.MotionValid == true) {
		dev := system.Device{
			Id:        s.EntityID,
			Namespace: namespace,
			Name:      s.Attributes.FriendlyName,
			Type:      system.MotionSensor,
		}
		return &dev

	} else if s.Attributes.DeviceClass == "opening" || s.Attributes.DeviceClass == "door" || s.Attributes.DeviceClass == "window" || s.Attributes.DeviceClass == "garage_door" {
		dev := system.Device{
			Id:        s.EntityID,
			Namespace: namespace,
			Name:      s.Attributes.FriendlyName,
			Type:      system.ContactSensor,
		}
		return &dev
	} else if s.Attributes.DeviceClass == "smoke" {
		dev := system.Device{
			Id:        s.EntityID,
			Namespace: namespace,
			Name:      s.Attributes.FriendlyName,
			Type:      system.SmokeSensor,
		}
		return &dev
	}
	return nil
}

func (ha *homeassistant) disconnectedHandler(controller dd.SystemController) connector.DisconnectedHandler {
	return func(err error) {
		if err != nil {
			log.Error().Err(err).Msg("Connection to HomeAssistant lost. Retrying in a few seconds...")
			ha.reconnect(controller)
		}
	}
}

func (ha *homeassistant) reconnect(controller dd.SystemController) {
	go func() {
		ha.Stop()
		<-time.After(10 * time.Second)
		ha.Start(controller)
	}()
}

func (ha *homeassistant) setupDevice(dev system.Device, batteryEntityId string, tamperEntityId string, controller dd.SystemController) {
	var sId uint64 = 0
	var err error = nil

	switch dev.Type {
	case system.MotionSensor:
		sId, err = ha.connector.SubscribeStateTrigger(dev.Id, driver.MotionSensorHandler(&dev, controller))
	case system.ContactSensor:
		sId, err = ha.connector.SubscribeStateTrigger(dev.Id, driver.ContactSensorHandler(&dev, controller))
	case system.SmokeSensor:
		sId, err = ha.connector.SubscribeStateTrigger(dev.Id, driver.SmokeSensorHandler(&dev, controller))
	}

	li := system.DInfo(&dev)
	if err != nil {
		system.DError(&dev).Err(err).Msg("Could not setup HomeAssistant device!")
	} else {
		if len(batteryEntityId) > 0 {
			sId, err := ha.connector.SubscribeStateTrigger(batteryEntityId, driver.BatteryHandler(&dev, controller))
			if err != nil {
				log.Error().Str("_battery", batteryEntityId).Str("_id", dev.Id).Err(err).Msg("Could not setup HomeAssistant battery tracker!")
			}
			ha.devices.Store(dev, sId)
			li.Str("_battery", batteryEntityId)
		}
		if len(tamperEntityId) > 0 {
			sId, err := ha.connector.SubscribeStateTrigger(tamperEntityId, driver.TamperHandler(&dev, controller))
			if err != nil {
				log.Error().Str("_tamper", tamperEntityId).Str("_id", dev.Id).Err(err).Msg("Could not setup HomeAssistant tamper tracker!")
			}
			ha.devices.Store(dev, sId)
			li.Str("_tamper", tamperEntityId)
		}
	}
	ha.devices.Store(dev, sId)
	li.Msg("Setup HomeAssistant device")
}

func (ha *homeassistant) tearDownAllDevices(connectionLost bool) {
	ha.devices.Range(func(key, value any) bool {
		if !connectionLost {
			ha.connector.UnsubscribeStateTrigger(value.(uint64))
		}
		dev := key.(system.Device)
		system.DInfo(&dev).Msg("Remove HomeAssistant device")
		return true
	})
	ha.devices = sync.Map{}
}

func (ha *homeassistant) Stop() {
	ha.connection += 1
	ha.tearDownAllDevices(true)
	ha.connector.Disconnect()
	ha.connector = nil
}

func (ha *homeassistant) getStates() (msgs.StateResult, error) {
	var result msgs.StateResult
	err := ha.connector.Command(&msgs.BaseMessage{Type: msgs.GetStates}, &result)
	return result, err
}

func (ha *homeassistant) testConnection(conId uint64) {
	for ha.connection == conId {
		ha.connector.Command(&msgs.BaseMessage{Type: msgs.Ping}, nil)
		<-time.After(30 * time.Second)
	}

}

func (ha *homeassistant) getDeviceConfig() map[string]devicesConfig {
	configs := cfg.GetObjects[devicesConfig](cDevices)
	result := make(map[string]devicesConfig)
	for _, c := range configs {
		result[c.Id] = c
	}
	return result
}
