package homeassistant

import "github.com/mtrossbach/waechter/system"

type devicesConfig struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	Battery string            `json:"battery"`
	Type    system.DeviceType `json:"type"`
}
