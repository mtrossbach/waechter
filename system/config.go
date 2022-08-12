package system

import "github.com/mtrossbach/waechter/config"

const (
	cExitDelay            = "general.exitdelay"
	cEntryDelay           = "general.entrydelay"
	cTamperAlarm          = "general.tamperalarm"
	cMaxWrongPinCount     = "general.maxwrongpincount"
	cBatteryThreshold     = "general.batterythreshold"
	cLinkQualityThreshold = "general.linkqualitythreshold"
	cDisarmPins           = "disarmpins"
)

func init() {
	config.SetDefault(cExitDelay, 30)
	config.SetDefault(cEntryDelay, 30)
	config.SetDefault(cTamperAlarm, true)
	config.SetDefault(cMaxWrongPinCount, 3)
	config.SetDefault(cBatteryThreshold, 0.1)
	config.SetDefault(cLinkQualityThreshold, 0.1)
	config.SetDefault(cDisarmPins, []string{})
}
