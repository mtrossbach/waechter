package i18n

type Key string

const (
	WALowBattery     Key = "whatsapp_low_battery"
	WALowLinkQuality Key = "whatsapp_low_link_quality"

	AlarmNone    Key = "alarm_none"
	AlarmBurglar Key = "alarm_burglar"
	AlarmPanic   Key = "alarm_panic"
	AlarmFire    Key = "alarm_fire"
	AlarmTamper  Key = "alarm_tamper"
)