package connector

import (
	"github.com/mtrossbach/waechter/internal/cfg"
)

const (
	cConnection = "zigbee2mqtt.connection"
	cClientId   = "zigbee2mqtt.clientid"
	cUsername   = "zigbee2mqtt.username"
	cPassword   = "zigbee2mqtt.password"
	cBaseTopic  = "zigbee2mqtt.basetopic"
)

func init() {
	cfg.SetDefault(cConnection, "mqtt://localhost:1883")
	cfg.SetDefault(cBaseTopic, "zigbee2mqtt")
}
