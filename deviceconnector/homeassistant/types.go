package homeassistant

import (
	"github.com/mtrossbach/waechter/system/device"
	"strings"
)

type assembledDevice struct {
	sensors     map[device.Sensor]string
	spec        device.Spec
	entityId    string
	displayName string
}

func (a *assembledDevice) generateSpec(connectorId string) device.Spec {
	spec := device.Spec{
		Id:          device.NewId(connectorId, a.entityId),
		DisplayName: a.displayName,
		Sensors:     []device.Sensor{},
		Actors:      []device.Actor{},
	}

	for s := range a.sensors {
		spec.Sensors = append(spec.Sensors, s)
	}
	return spec
}

func entityPrefix(entityId string) string {
	if strings.Contains(entityId, "_") {
		return entityId[:strings.LastIndex(entityId, "_")] + "_"
	}
	return entityId
}
