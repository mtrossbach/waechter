package zigbee2mqtt

import (
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/cfg"
)

const (
	cConnection          = "zigbee2mqtt.connection"
	cClientId            = "zigbee2mqtt.clientid"
	cUsername            = "zigbee2mqtt.username"
	cPassword            = "zigbee2mqtt.password"
	cBaseTopic           = "zigbee2mqtt.basetopic"
	cEnable              = "zigbee2mqtt.enable"
	cAutoDeviceDiscovery = "zigbee2mqtt.autodevicediscovery"
	cDevices             = "zigbee2mqtt.devices"
)

func init() {
	cfg.SetDefault(cEnable, false)
	cfg.SetDefault(cConnection, "mqtt://localhost:1883")
	cfg.SetDefault(cBaseTopic, "zigbee2mqtt")
	cfg.SetDefault(cAutoDeviceDiscovery, true)
}

func cOptions() connector.Options {
	return connector.Options{
		Uri:       cfg.GetString(cConnection),
		ClientId:  cfg.GetString(cClientId),
		Username:  cfg.GetString(cUsername),
		Password:  cfg.GetString(cPassword),
		BaseTopic: cfg.GetString(cBaseTopic),
	}
}

func IsEnabled() bool {
	return cfg.GetBool(cEnable)
}
