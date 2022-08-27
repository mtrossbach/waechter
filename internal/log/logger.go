package log

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"os"
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
