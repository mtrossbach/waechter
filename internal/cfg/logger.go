package cfg

import (
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

func init() {
	SetDefault(cLogLevel, "info")
	SetDefault(cLogFormat, "text")
}

func updateLogger() {
	level := strings.ToLower(GetString(cLogLevel))

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

	format := strings.ToLower(GetString(cLogFormat))
	switch format {
	case "json":
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	default:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
}

func Logger(component string) zerolog.Logger {
	return log.With().Str("_comp", component).Logger()
}
