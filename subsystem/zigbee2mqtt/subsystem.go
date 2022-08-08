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
)

type Zigbee2MqttSubsystem struct {
	deviceManager system.DeviceManager
	z2mManager    *zigbee.Z2MManager
}

func NewZigbee2MqttSubsystem() *Zigbee2MqttSubsystem {
	return &Zigbee2MqttSubsystem{
		z2mManager: zigbee.NewZ2MManager(config.GetConfig().Zigbee2Mqtt),
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
		misc.Log.Warnf("Could not parse device event: %v", string(msg.Payload()))
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
		misc.Log.Warnf("Could not parse devices payload: %v", string(msg.Payload()))
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

	misc.Log.Debugf("Received new list of device ids: %v", relevantDeviceIds)
	misc.Log.Debugf("Already registered device ids: %v", oldDeviceIds)
	misc.Log.Debugf("Device ids to add: %v", deviceIdsToAdd)
	misc.Log.Debugf("Device ids to remove: %v", deviceIdsToRemove)

	for _, removeId := range deviceIdsToRemove {
		z2ms.deviceManager.RemoveDeviceById(removeId)
	}

	for _, addId := range deviceIdsToAdd {
		dev := devices.CreateDevice(relevantDevices[addId], z2ms.z2mManager)
		if dev != nil {
			z2ms.deviceManager.AddDevice(dev)
		} else {
			misc.Log.Warnf("Could not find driver for device with id %v: %v", addId, relevantDevices[addId])
		}
	}
}
