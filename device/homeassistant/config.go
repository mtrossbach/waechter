package homeassistant

import "github.com/mtrossbach/waechter/internal/cfg"

const (
	cURL                 = "homeassistant.url"
	cToken               = "homeassistant.token"
	cEnable              = "homeassistant.enable"
	cAutoDeviceDiscovery = "homeassistant.autodevicediscovery"
	cDevices             = "homeassistant.devices"
)

func init() {
	cfg.SetDefault(cAutoDeviceDiscovery, true)
}

func IsEnabled() bool {
	return cfg.GetBool(cEnable)
}
