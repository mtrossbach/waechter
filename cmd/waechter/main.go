package main

import (
	"github.com/mtrossbach/waechter/device/homeassistant"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/internal/i18n"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/notification/dummy"
	"github.com/mtrossbach/waechter/notification/whatsapp"
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
	i18n.InitI18n()

	sys := system.NewWaechterSystem()
	go zigbee2mqtt.New().Start(sys)
	go homeassistant.New().Start(sys)
	sys.AddNotificationAdapter(whatsapp.NewWhatsApp())
	sys.AddNotificationAdapter(dummy.New())

	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Debug().Str("sig", sig.String()).Msg("Caught SIGTERM")
	log.Info().Msg("Exit.")
}
