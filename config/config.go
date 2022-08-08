package config

import (
	"sync"
	"time"
)

var once sync.Once
var instance *Config

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			General: General{
				ExitDelay:            30,
				EntryDelay:           30,
				TamperAlarm:          true,
				Siren:                true,
				MaxWrongPinCount:     3,
				BatteryThresold:      0.2,
				LinkQualityThreshold: 0.2,
			},
			Zigbee2Mqtt: &Zigbee2Mqtt{
				Connection: "mqtt://waechterpi:1883",
				ClientId:   "",
				Username:   "",
				Password:   "",
				BaseTopic:  "zigbee2mqtt",
			},
			DisarmPins: []DisarmPin{
				{
					Pin:         "1337",
					DisplayName: "Markus",
				},
			},
		}
	})
	return instance
}

type Config struct {
	General     General
	Zigbee2Mqtt *Zigbee2Mqtt
	DisarmPins  []DisarmPin
}

type General struct {
	ExitDelay            time.Duration
	EntryDelay           time.Duration
	TamperAlarm          bool
	Siren                bool
	MaxWrongPinCount     int
	BatteryThresold      float32
	LinkQualityThreshold float32
}

type Zigbee2Mqtt struct {
	Connection string
	ClientId   string
	Username   string
	Password   string
	BaseTopic  string
}

type DisarmPin struct {
	Pin         string
	DisplayName string
}
