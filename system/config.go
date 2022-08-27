package system

import (
	"github.com/mtrossbach/waechter/internal/cfg"
)

const (
	cExitDelay            = "general.exitdelay"
	cEntryDelay           = "general.entrydelay"
	cTamperAlarm          = "general.tamperalarm"
	cMaxWrongPinCount     = "general.maxwrongpincount"
	cBatteryThreshold     = "general.batterythreshold"
	cLinkQualityThreshold = "general.linkqualitythreshold"
	cDisarmPins           = "disarmpins"

	cSystemState      = "system.state"
	cSystemArmingMode = "system.armingMode"
	cSystemAlarmType  = "system.alarmType"
)

func init() {
	cfg.SetDefault(cExitDelay, 30)
	cfg.SetDefault(cEntryDelay, 30)
	cfg.SetDefault(cTamperAlarm, true)
	cfg.SetDefault(cMaxWrongPinCount, 3)
	cfg.SetDefault(cBatteryThreshold, 0.1)
	cfg.SetDefault(cLinkQualityThreshold, 0.05)
	cfg.SetDefault(cDisarmPins, []string{})
}
