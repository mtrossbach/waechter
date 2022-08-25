package zigbee2mqtt

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/internal/wslice"

	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/system"
)

type zigbee2mqtt struct {
	deviceManager system.DeviceSystem
	connector     *connector.Connector
}

func New() *zigbee2mqtt {

	return &zigbee2mqtt{
		connector: connector.New(),
	}
}

func (z2ms *zigbee2mqtt) GetName() string {
	return model2.SubsystemName
}

func (z2ms *zigbee2mqtt) Start(deviceManager system.DeviceSystem) {
	z2ms.deviceManager = deviceManager
	z2ms.connector.Connect()
	z2ms.connector.Subscribe("bridge/devices", z2ms.handleNewDeviceList)
	z2ms.connector.Subscribe("bridge/event", z2ms.handleDeviceEvent)
}

func (z2ms *zigbee2mqtt) Stop() {
	z2ms.connector.Disconnect()
}

func (z2ms *zigbee2mqtt) handleDeviceEvent(msg mqtt.Message) {
	var deviceEvent model2.DeviceEvent
	if err := json.Unmarshal(msg.Payload(), &deviceEvent); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse zdevice event!")
		return
	}

	if deviceEvent.Type == "device_announce" && len(deviceEvent.Data.IeeeAddress) > 0 {
		device := z2ms.deviceManager.GetDeviceById(deviceEvent.Data.IeeeAddress)
		zdev, ok := device.(zdevice.ZDevice)
		if ok {
			zdev.OnDeviceAnnounced()
		}
	}
}

func (z2ms *zigbee2mqtt) handleNewDeviceList(msg mqtt.Message) {
	var newDevices []model2.Z2MDeviceInfo
	if err := json.Unmarshal(msg.Payload(), &newDevices); err != nil {
		log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse devices payload!")
		return
	}

	var relevantDeviceIds []string
	relevantDevices := make(map[string]model2.Z2MDeviceInfo)
	for _, device := range newDevices {
		if device.Type == "EndDevice" && device.Supported {
			relevantDeviceIds = append(relevantDeviceIds, device.IeeeAddress)
			relevantDevices[device.IeeeAddress] = device
		}
	}

	oldDeviceIds := z2ms.deviceManager.GetDeviceIdsForSubsystem(z2ms.GetName())

	deviceIdsToRemove := wslice.StringsMissingInList(relevantDeviceIds, oldDeviceIds)
	deviceIdsToAdd := wslice.StringsMissingInList(oldDeviceIds, relevantDeviceIds)

	log.Debug().Strs("new", relevantDeviceIds).Strs("old", oldDeviceIds).Msg("Received new device list")
	log.Debug().Strs("add", deviceIdsToAdd).Strs("remove", deviceIdsToRemove).Msg("Merging device list")

	for _, removeId := range deviceIdsToRemove {
		z2ms.deviceManager.RemoveDeviceById(removeId)
	}

	for _, addId := range deviceIdsToAdd {
		dev := zdevice.CreateDevice(relevantDevices[addId], z2ms.connector)
		if dev != nil {
			z2ms.deviceManager.AddDevice(dev)
		} else {
			log.Error().Str("id", addId).Interface("zdevice", relevantDevices[addId]).Msg("Could not find driver for device. Device will not be added to the system.")
		}
	}
}
