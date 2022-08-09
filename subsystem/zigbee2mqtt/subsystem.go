package zigbee2mqtt

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/devices"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type Zigbee2MqttSubsystem struct {
	deviceManager system.DeviceManager
	z2mManager    *zigbee.Z2MManager
	log           zerolog.Logger
}

func NewZigbee2MqttSubsystem() *Zigbee2MqttSubsystem {
	return &Zigbee2MqttSubsystem{
		z2mManager: zigbee.NewZ2MManager(config.GetConfig().Zigbee2Mqtt),
		log:        misc.Logger("Zigbee2MqttSubsystem"),
	}
}

func (z2ms *Zigbee2MqttSubsystem) GetName() string {
	return model.SubsystemName
}

func (z2ms *Zigbee2MqttSubsystem) Start(deviceManager system.DeviceManager) {
	z2ms.deviceManager = deviceManager
	z2ms.z2mManager.Connect()
	z2ms.z2mManager.Subscribe("bridge/devices", z2ms.handleNewDeviceList)
	z2ms.z2mManager.Subscribe("bridge/events", z2ms.handleDeviceEvent)
}

func (z2ms *Zigbee2MqttSubsystem) Stop() {
	z2ms.z2mManager.Disconnect()
}

func (z2ms *Zigbee2MqttSubsystem) handleDeviceEvent(msg mqtt.Message) {
	var deviceEvent model.DeviceEvent
	if err := json.Unmarshal(msg.Payload(), &deviceEvent); err != nil {
		z2ms.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse device event!")
		return
	}

	if deviceEvent.Type == "device_announce" && len(deviceEvent.Data.IeeeAddress) > 0 {
		device := z2ms.deviceManager.GetDeviceById(deviceEvent.Data.IeeeAddress)
		zdevice, ok := device.(devices.ZDevice)
		if ok {
			zdevice.OnDeviceAnnounced()
		}
	}
}

func (z2ms *Zigbee2MqttSubsystem) handleNewDeviceList(msg mqtt.Message) {
	var newDevices []model.Z2MDeviceInfo
	if err := json.Unmarshal(msg.Payload(), &newDevices); err != nil {
		z2ms.log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse devices payload!")
		return
	}

	var relevantDeviceIds []string
	relevantDevices := make(map[string]model.Z2MDeviceInfo)
	for _, device := range newDevices {
		if device.Type == "EndDevice" && device.Supported {
			relevantDeviceIds = append(relevantDeviceIds, device.IeeeAddress)
			relevantDevices[device.IeeeAddress] = device
		}
	}

	oldDeviceIds := z2ms.deviceManager.GetDeviceIdsForSubsystem(z2ms.GetName())

	deviceIdsToRemove := misc.StringsMissingInList(relevantDeviceIds, oldDeviceIds)
	deviceIdsToAdd := misc.StringsMissingInList(oldDeviceIds, relevantDeviceIds)

	z2ms.log.Debug().Strs("new", relevantDeviceIds).Strs("old", oldDeviceIds).Strs("add", deviceIdsToAdd).Strs("remove", deviceIdsToRemove).Msg("Processing device list")

	for _, removeId := range deviceIdsToRemove {
		z2ms.deviceManager.RemoveDeviceById(removeId)
	}

	for _, addId := range deviceIdsToAdd {
		dev := devices.CreateDevice(relevantDevices[addId], z2ms.z2mManager)
		if dev != nil {
			z2ms.deviceManager.AddDevice(dev)
		} else {
			z2ms.log.Warn().Str("id", addId).Interface("device", relevantDevices[addId]).Msg("Could not finde driver for device")
		}
	}
}
