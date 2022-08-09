package misc

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once sync.Once

func Logger(component string) zerolog.Logger {
	once.Do(func() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	})
	return log.With().Str("component", component).Logger()
}
