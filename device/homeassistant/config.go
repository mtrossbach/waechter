package homeassistant

import "github.com/mtrossbach/waechter/internal/cfg"

const (
	cURL    = "homeassistant.url"
	cToken  = "homeassistant.token"
	cEnable = "homeassistant.enable"
)

func IsEnabled() bool {
	return cfg.GetBool(cEnable)
}
