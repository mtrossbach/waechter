package zigbee2mqtt

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/connector"
	model2 "github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt/zdevice"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type zigbee2mqtt struct {
	deviceManager system.DeviceSystem
	z2mManager    *connector.Z2MManager
	log           zerolog.Logger
}

func New() *zigbee2mqtt {
	return &zigbee2mqtt{
		z2mManager: connector.NewZ2MManager(config.GetConfig().Zigbee2Mqtt),
		log:        misc.Logger("Zigbee2Mqtt"),
	}
}

func (z2ms *zigbee2mqtt) GetName() string {
	return model2.SubsystemName
}

func (z2ms *zigbee2mqtt) Start(deviceManager system.DeviceSystem) {
	z2ms.deviceManager = deviceManager
	z2ms.z2mManager.Connect()
	z2ms.z2mManager.Subscribe("bridge/devices", z2ms.handleNewDeviceList)
	z2ms.z2mManager.Subscribe("bridge/events", z2ms.handleDeviceEvent)
}

func (z2ms *zigbee2mqtt) Stop() {
	z2ms.z2mManager.Disconnect()
}

func (z2ms *zigbee2mqtt) handleDeviceEvent(msg mqtt.Message) {
	var deviceEvent model2.DeviceEvent
	if err := json.Unmarshal(msg.Payload(), &deviceEvent); err != nil {
		z2ms.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse zdevice event!")
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
		z2ms.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse devices payload!")
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

	deviceIdsToRemove := misc.StringsMissingInList(relevantDeviceIds, oldDeviceIds)
	deviceIdsToAdd := misc.StringsMissingInList(oldDeviceIds, relevantDeviceIds)

	z2ms.log.Debug().Strs("new", relevantDeviceIds).Strs("old", oldDeviceIds).Msg("Received new device list")
	z2ms.log.Debug().Strs("add", deviceIdsToAdd).Strs("remove", deviceIdsToRemove).Msg("Merging device list")

	for _, removeId := range deviceIdsToRemove {
		z2ms.deviceManager.RemoveDeviceById(removeId)
	}

	for _, addId := range deviceIdsToAdd {
		dev := zdevice.CreateDevice(relevantDevices[addId], z2ms.z2mManager)
		if dev != nil {
			z2ms.deviceManager.AddDevice(dev)
		} else {
			z2ms.log.Warn().Str("id", addId).Interface("zdevice", relevantDevices[addId]).Msg("Could not find driver for device. Device will not be added to the system.")
		}
	}
}
