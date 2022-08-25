package main

import (
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant"
	"github.com/mtrossbach/waechter/subsystem/device/zigbee2mqtt"
	"github.com/mtrossbach/waechter/subsystem/notif/dummy"
	"github.com/mtrossbach/waechter/system"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg.Init()
	cfg.Print()
	log.UpdateLogger()

	log.Info().Msg("Starting up...")
	sys := system.NewWaechterSystem()
	zigbee2mqtt.New().Start(sys)
	homeassistant.New().Start(sys)
	sys.AddNotificationHandler(dummy.New().SendNotification)
	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Debug().Str("sig", sig.String()).Msg("Caught SIGTERM")
	log.Info().Msg("Exit.")
}
