package config

import (
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

var file string
var instance *Config

func Init() {
	possiblePaths := []string{
		"./config.yaml",
		"./config/config.yaml",
		"/config.yaml",
		"~/waechter/config.yaml",
		"~/.waechter/config.yaml",
		"/etc/waechter/config.yaml",
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			file, _ = filepath.Abs(p)
			break
		}
	}

	if len(file) == 0 {
		log.Fatalf("No config file found!\n")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Could not open config: %v", err)
	}

	var config Config
	defaults.Set(&config)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Could not read config: %v", err)
	}

	instance = &config
}

func General() GeneralConfig {
	return instance.General
}

func Log() LogConfig {
	return instance.Log
}

func Persons() []Person {
	return instance.Persons
}

func Devices() []DeviceConfig {
	return instance.Devices
}

func Zones() []ZoneConfig {
	return instance.Zones
}

func Zigbee2Mqtt() []Zigbee2MqttConfig {
	return instance.Zigbee2Mqtt
}

func HomeAssistant() []HomeAssistantConfig {
	return instance.HomeAssistant
}

func WhatsApp() *WhatsAppConfiguration {
	return instance.WhatsApp
}

func Notification() []string {
	return instance.Notification
}
