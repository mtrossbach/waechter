package homeassistant

import (
	dd "github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/homeassistant/connector"
	"github.com/mtrossbach/waechter/device/homeassistant/driver"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
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
	log.Info().Str("uri", cfg.GetString(cURL)).Msg("Connecting to HomeAssistant...")
	err := ha.connector.Connect(cfg.GetString(cURL), cfg.GetString(cToken), ha.disconnectedHandler(controller))
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to HomeAssistant. Retrying in a few seconds...")
		ha.reconnect(controller)
		return
	}
	go ha.testConnection(ha.connection)
	st, err := ha.getStates()
	if err != nil {
		log.Error().Err(err).Msg("Could not request states from HomeAssistant")
		return
	}

	var devices []system.Device

	for _, s := range st.Result {
		if s.Attributes.DeviceClass == "motion" && (s.Attributes.MotionValid == nil || *s.Attributes.MotionValid == true) {
			dev := system.Device{
				Id:        s.EntityID,
				Namespace: namespace,
				Name:      s.Attributes.FriendlyName,
				Type:      system.MotionSensor,
			}
			devices = append(devices, dev)

		} else if s.Attributes.DeviceClass == "opening" || s.Attributes.DeviceClass == "door" || s.Attributes.DeviceClass == "window" || s.Attributes.DeviceClass == "garage_door" {
			dev := system.Device{
				Id:        s.EntityID,
				Namespace: namespace,
				Name:      s.Attributes.FriendlyName,
				Type:      system.ContactSensor,
			}
			devices = append(devices, dev)
		} else if s.Attributes.DeviceClass == "smoke" {
			dev := system.Device{
				Id:        s.EntityID,
				Namespace: namespace,
				Name:      s.Attributes.FriendlyName,
				Type:      system.SmokeSensor,
			}
			devices = append(devices, dev)
		}
	}

	for _, dev := range devices {
		prefix := strings.Replace(dev.Id[:strings.LastIndex(dev.Id, "_")], "binary_sensor.", "sensor.", 1)
		var batteryEntityId string
		for _, s := range st.Result {
			if strings.HasPrefix(s.EntityID, prefix) && s.Attributes.DeviceClass == "battery" {
				batteryEntityId = s.EntityID
				break
			}
		}
		ha.setupDevice(dev, batteryEntityId, controller)
	}

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

func (ha *homeassistant) setupDevice(dev system.Device, batteryEntityId string, controller dd.SystemController) {
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

	if err != nil {
		system.DError(&dev).Err(err).Msg("Could not setup HomeAssistant device!")
	} else {
		system.DInfo(&dev).Msg("Setup HomeAssistant device")
	}
	ha.devices.Store(dev, sId)

	if len(batteryEntityId) > 0 {
		sId, err := ha.connector.SubscribeStateTrigger(batteryEntityId, driver.BatteryHandler(&dev, controller))
		if err != nil {
			log.Error().Str("_batteryId", batteryEntityId).Str("_id", dev.Id).Err(err).Msg("Could not setup HomeAssistant battery tracker!")
		} else {
			log.Info().Str("_batteryId", batteryEntityId).Str("_id", dev.Id).Msg("Setup HomeAssistant battery tracker")
		}
		ha.devices.Store(dev, sId)
	}
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
