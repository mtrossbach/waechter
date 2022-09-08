package whatsapp

import "github.com/mtrossbach/waechter/internal/cfg"

const (
	cEnable  = "whatsapp.enable"
	cToken   = "whatsapp.token"
	cPhoneId = "whatsapp.phoneid"

	cAlarmTemplateName        = "whatsapp.template.alarm"
	cRecoverTemplateName      = "whatsapp.template.recover"
	cNotificationTemplateName = "whatsapp.template.notification"
	cAutoArmTemplateName      = "whatsapp.template.autoarm"
	cAutoDisarmTemplateName   = "whatsapp.template.autodisarm"
)

func IsEnabled() bool {
	return cfg.GetBool(cEnable)
}
