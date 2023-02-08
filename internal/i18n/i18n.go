package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system/alarm"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path"
	"path/filepath"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	possiblePaths := []string{
		"./locales",
		"/locales",
		"~/waechter/locales",
		"~/.waechter/locales",
		"/etc/waechter/locales",
	}

	basePath := ""
	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			basePath, _ = filepath.Abs(p)
			break
		}
	}

	if len(basePath) > 0 {
		log.Info().Str("path", basePath).Msg("Found localizations.")
	}

	_, err := bundle.LoadMessageFile(path.Join(basePath, "en.json"))
	if err != nil {
		log.Error().Err(err).Msg("Could not load en.json localization file.")
	}
	_, err = bundle.LoadMessageFile(path.Join(basePath, "de.json"))
	if err != nil {
		log.Error().Err(err).Msg("Could not load de.json localization file.")
	}
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
