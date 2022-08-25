package log

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	cLogLevel  = "log.level"
	cLogFormat = "log.format"
)

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func TDebug(t any) *zerolog.Event {
	if t == nil {
		return Debug()
	}
	return log.Debug().Str("_c", reflect.TypeOf(t).String())
}

func TInfo(t any) *zerolog.Event {
	if t == nil {
		return Info()
	}
	return log.Info().Str("_c", reflect.TypeOf(t).String())
}

func TError(t any) *zerolog.Event {
	if t == nil {
		return Error()
	}
	return log.Error().Str("_c", reflect.TypeOf(t).String())
}

func init() {
	cfg.SetDefault(cLogLevel, "info")
	cfg.SetDefault(cLogFormat, "text")
}

func UpdateLogger() {
	level := strings.ToLower(cfg.GetString(cLogLevel))

	switch level {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	format := strings.ToLower(cfg.GetString(cLogFormat))
	switch format {
	case "json":
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	default:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
}
