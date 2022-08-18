package siren

import (
	"github.com/mtrossbach/waechter/internal/cfg"
)

const (
	cEnabled = "siren.enabled"
	cLevel   = "siren.level"
)

func init() {
	cfg.SetDefault(cEnabled, true)
	cfg.SetDefault(cLevel, "high")
}
