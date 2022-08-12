package siren

import "github.com/mtrossbach/waechter/config"

const (
	cEnabled = "siren.enabled"
	cLevel   = "siren.level"
)

func init() {
	config.SetDefault(cEnabled, true)
	config.SetDefault(cLevel, "high")
}
