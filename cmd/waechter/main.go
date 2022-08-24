package main

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant"
	"github.com/mtrossbach/waechter/subsystem/notif/dummy"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg.Init()
	cfg.Print()

	log.Info().Msg("Starting up...")
	sys := system.NewWaechterSystem()
	//sys.RegisterDeviceSubsystem(zigbee2mqtt.New())
	sys.RegisterDeviceSubsystem(homeassistant.New())
	sys.RegisterNotifSubsystem(dummy.New())
	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Debug().Str("sig", sig.String()).Msg("Caught SIGTERM")
	log.Info().Msg("Exit.")
}
