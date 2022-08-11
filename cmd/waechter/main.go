package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mtrossbach/waechter/config"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt"
	"github.com/mtrossbach/waechter/subsystem/notif/dummy"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config.Init()
	config.Print()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Msg("Starting up...")
	sys := system.NewWaechterSystem()
	sys.RegisterDeviceSubsystem(zigbee2mqtt.New())
	sys.RegisterNotifSubsystem(dummy.New())
	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Debug().Str("sig", sig.String()).Msg("Caught SIGTERM")
	log.Info().Msg("Exit.")
}
