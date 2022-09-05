package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
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

func TranslateAlarm(lang string, alarmType system.AlarmType) string {
	switch alarmType {
	case system.NoAlarm:
		return Translate(lang, AlarmNone)
	case system.BurglarAlarm:
		return Translate(lang, AlarmBurglar)
	case system.PanicAlarm:
		return Translate(lang, AlarmPanic)
	case system.FireAlarm:
		return Translate(lang, AlarmFire)
	case system.TamperAlarm:
		return Translate(lang, AlarmTamper)
	}

	return string(alarmType)
}
