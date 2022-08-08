package devices

import (
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
)

type ZDevice interface {
	system.Device
	OnDeviceAnnounced()
}

func CreateDevice(deviceInfo model.Z2MDeviceInfo, z2mManager *zigbee.Z2MManager) system.Device {
	var exposes []string

	for _, e := range deviceInfo.Definition.Exposes {
		exposes = append(exposes, e.Property)
	}

	if misc.ContainsAll(exposes, []string{"action_code", "action"}) {
		return newGenericKeypad(deviceInfo, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"alarm"}) {
		return newGenericSiren(deviceInfo, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"contact"}) {
		return newGenericContactSensor(deviceInfo, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"occupancy"}) {
		return newGenericMotionSensor(deviceInfo, z2mManager)
	} else {
		return nil
	}
}
