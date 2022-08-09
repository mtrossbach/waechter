package runner

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt"
	"github.com/mtrossbach/waechter/system"
)

func Run() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting up...")

	system := system.NewWaechterSystem()
	system.RegisterSubsystem(zigbee2mqtt.NewZigbee2MqttSubsystem())
	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Info().Str("sig", sig.String()).Msg("Caught SIGTERM")
}
