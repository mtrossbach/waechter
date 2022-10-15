package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("./locales/en.json")
	bundle.LoadMessageFile("./locales/de.json")
}

func Translate(lang string, key Key) string {
	localizer := i18n.NewLocalizer(bundle, lang)

	val, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: string(key),
	})

	if err != nil {
		log.Error().Err(err).Str("key", string(key)).Msg("Could not translate message")
		return fmt.Sprintf("###%v###", key)
	}
	return val
}

func TranslateAlarm(lang string, alarmType alarm.Type) string {
	switch alarmType {
	case alarm.None:
		return Translate(lang, AlarmNone)
	case alarm.EntryDelay:
		return Translate(lang, AlarmEntryDelay)
	case alarm.Burglar:
		return Translate(lang, AlarmBurglar)
	case alarm.Panic:
		return Translate(lang, AlarmPanic)
	case alarm.Fire:
		return Translate(lang, AlarmFire)
	case alarm.Tamper:
		return Translate(lang, AlarmTamper)
	case alarm.TamperPin:
		return Translate(lang, AlarmTamperPin)
	}
	return string(alarmType)
}
