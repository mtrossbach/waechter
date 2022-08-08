package devices

import (
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/model"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt/zigbee"
	"github.com/mtrossbach/waechter/system"
)

func CreateDevice(zdevice model.ZigbeeDevice, z2mManager *zigbee.Z2MManager) system.Device {
	var exposes []string

	for _, e := range zdevice.Definition.Exposes {
		exposes = append(exposes, e.Property)
	}

	if misc.ContainsAll(exposes, []string{"action_code", "action"}) {
		return newGenericKeypad(zdevice, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"alarm"}) {
		return newGenericSiren(zdevice, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"contact"}) {
		return newGenericContactSensor(zdevice, z2mManager)
	} else if misc.ContainsAll(exposes, []string{"occupancy"}) {
		return newGenericMotionSensor(zdevice, z2mManager)
	} else {
		return nil
	}
}
