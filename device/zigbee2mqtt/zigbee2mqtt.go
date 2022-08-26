package zigbee2mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/system"
)

type zigbee2mqtt struct {
	systemController device.SystemController
	connector        *connector.Connector
	devices          sync.Map
}

func New() *zigbee2mqtt {
	return &zigbee2mqtt{
		connector: connector.New(),
	}
}

func (z2ms *zigbee2mqtt) Start(systemController device.SystemController) {
	z2ms.systemController = systemController
	systemController.SubscribeStateUpdate(z2ms, z2ms.updateState)
	z2ms.connector.Connect()
	z2ms.connector.Subscribe("bridge/devices", z2ms.handleNewDeviceList)
	z2ms.connector.Subscribe("bridge/event", z2ms.handleDeviceEvent)
}

func (z2ms *zigbee2mqtt) updateState(state system.State, armingMode system.ArmingMode, alarmType system.AlarmType) {
	z2ms.devices.Range(func(_, value any) bool {
		(value.(ZDevice)).UpdateState(state, armingMode, alarmType)
		return true
	})
}

func (z2ms *zigbee2mqtt) Stop() {
	z2ms.connector.Disconnect()
}

func (z2ms *zigbee2mqtt) handleDeviceEvent(msg mqtt.Message) {
	var deviceEvent DeviceEvent
	if err := json.Unmarshal(msg.Payload(), &deviceEvent); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse zdevice event!")
		return
	}

	if deviceEvent.Type == "device_announce" && len(deviceEvent.Data.IeeeAddress) > 0 {
		dev, ok := z2ms.devices.Load(ieee2Id(deviceEvent.Data.IeeeAddress))
		if ok {
			zdev, ok := dev.(ZDevice)
			if ok {
				zdev.OnDeviceAnnounced()
			}
		}
	}
}

func (z2ms *zigbee2mqtt) handleNewDeviceList(msg mqtt.Message) {
	var newDevices []Z2MDeviceInfo
	if err := json.Unmarshal(msg.Payload(), &newDevices); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse devices payload!")
		return
	}

	relevantDevices := make(map[string]Z2MDeviceInfo)
	for _, device := range newDevices {
		if device.Type == "EndDevice" && device.Supported {
			relevantDevices[device.IeeeAddress] = device
		}
	}

	z2ms.devices.Range(func(_, value any) bool {
		value.(ZDevice).Teardown()
		return true
	})

	z2ms.devices = sync.Map{}
	for _, d := range relevantDevices {
		dev := createDevice(d)
		if dev != nil {
			dev.Setup(z2ms.connector, z2ms.systemController)
			z2ms.devices.Store(ieee2Id(d.IeeeAddress), dev)
		}
	}
}

func ieee2Id(ieeeAddress string) string {
	return fmt.Sprintf("%v-%v", "z2m", ieeeAddress)
}
