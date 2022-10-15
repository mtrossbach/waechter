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
	vendor      string
	model       string
	description string
}

func (a *assembledDevice) generateSpec(connectorId string) device.Spec {
	spec := device.Spec{
		Id:          device.NewId(connectorId, a.entityId),
		DisplayName: a.displayName,
		Vendor:      a.vendor,
		Model:       a.model,
		Description: a.description,
		Sensors:     []device.Sensor{},
		Actors:      []device.Actor{},
	}

	for s := range a.sensors {
		spec.Sensors = append(spec.Sensors, s)
	}
	return spec
}

func entityPrefix(entityId string) string {
	return entityId[:strings.LastIndex(entityId, "_")] + "_"
}
