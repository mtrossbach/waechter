package homeassistant

import (
	"github.com/mtrossbach/waechter/system"
	"strings"
)

type devicesConfig struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Battery string            `json:"battery"`
	Tamper  string            `json:"tamper"`
	Type    system.DeviceType `json:"type"`
}

func entityPrefix(entityId string) string {
	return entityId[:strings.LastIndex(entityId, "_")] + "_"
}
