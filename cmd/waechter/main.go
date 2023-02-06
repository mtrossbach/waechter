package main

import (
	"fmt"
	"github.com/mtrossbach/waechter/deviceconnector/homeassistant"
	"github.com/mtrossbach/waechter/deviceconnector/zigbee2mqtt"
	"github.com/mtrossbach/waechter/internal/config"
	"github.com/mtrossbach/waechter/internal/i18n"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/notification/whatsapp"
	"github.com/mtrossbach/waechter/system"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init()
	fmt.Printf("Using config file: %v\n", config.ConfigFile())
	config.Print()
	log.UpdateLogger()

	log.Info().Str("version", os.Getenv("WAECHTER_VERSION")).Msg("Starting up...")
	i18n.InitI18n()

	waechter := system.NewWaechter()
	z2ms := config.Zigbee2MqttConfigs()
	for _, z := range z2ms {
		if c, err := zigbee2mqtt.NewConnector(z); err != nil {
			log.Error().Err(err).Str("connector", "Zigbee2Mqtt").Str("id", z.Id).Msg("Could not initialize connector.")
		} else {
			waechter.AddDeviceConnector(c)
		}
	}
	has := config.HomeAssistantConfigs()
	for _, h := range has {
		if c, err := homeassistant.NewConnector(h); err != nil {
			log.Error().Err(err).Str("connector", "HomeAssistant").Str("id", h.Id).Msg("Could not initialize connector.")
		} else {
			waechter.AddDeviceConnector(c)
		}
	}

	for _, n := range config.Notifications() {
		switch n {
		case "whatsapp":
			if config := config.WhatsAppConfig(); config != nil {
				waechter.AddNotificationAdapter(whatsapp.NewWhatsApp(*config))
			}
		}
	}

	log.Info().Msg("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Debug().Str("sig", sig.String()).Msg("Caught SIGTERM")
	log.Info().Msg("Exit.")
}
