package connector

import (
	"github.com/mtrossbach/waechter/config"
)

const (
	cConnection = "zigbee2mqtt.connection"
	cClientId   = "zigbee2mqtt.clientid"
	cUsername   = "zigbee2mqtt.username"
	cPassword   = "zigbee2mqtt.password"
	cBaseTopic  = "zigbee2mqtt.basetopic"
)

func init() {
	config.SetDefault(cConnection, "mqtt://localhost:1883")
	config.SetDefault(cBaseTopic, "zigbee2mqtt")
}
