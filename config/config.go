package config

import (
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/mtrossbach/waechter/misc"
)

var once sync.Once
var instance *Config

func GetConfig() *Config {
	once.Do(func() {
		log := misc.Logger("Config")
		err := godotenv.Load()
		if err != nil {
			log.Info().Msg("No .env file present.")
		}

		instance = &Config{
			General: General{
				ExitDelay:            s2seconds(loadEnv("WAECHTER_GENERAL_EXITDELAY"), 0),
				EntryDelay:           s2seconds(loadEnv("WAECHTER_GENERAL_ENTRYDELAY"), 0),
				TamperAlarm:          s2bool(loadEnv("WAECHTER_GENERAL_TAMPERALARM"), true),
				Siren:                s2bool(loadEnv("WAECHTER_GENERAL_SIREN"), true),
				MaxWrongPinCount:     s2int(loadEnv("WAECHTER_GENERAL_WRONG_PIN_COUNT"), 3),
				BatteryThresold:      s2float32(loadEnv("WAECHTER_GENERAL_BATTERY_THRESHOLD"), 0.2),
				LinkQualityThreshold: s2float32(loadEnv("WAECHTER_GENERAL_LINKQUALITY_THREASHOLD"), 0.2),
			},
			Zigbee2Mqtt: &Zigbee2Mqtt{
				Connection: loadEnv("WAECHTER_ZIGBEE2MQTT_CONNECTION"),
				ClientId:   loadEnv("WAECHTER_ZIGBEE2MQTT_CLIENTID"),
				Username:   loadEnv("WAECHTER_ZIGBEE2MQTT_USERNAME"),
				Password:   loadEnv("WAECHTER_ZIGBEE2MQTT_PASSWORD"),
				BaseTopic:  loadEnv("WAECHTER_ZIGBEE2MQTT_BASETOPIC"),
			},
			DisarmPins: parsePin(loadEnv("WAECHTER_DISARMPINS")),
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

func loadEnv(key string) string {
	return os.Getenv(key)
}

func parsePin(input string) []DisarmPin {
	var result []DisarmPin
	items := strings.Split(input, ";")
	for _, item := range items {
		parts := strings.Split(item, ":")
		pin := DisarmPin{
			Pin: parts[0],
		}
		if len(parts) > 1 {
			pin.DisplayName = parts[1]
		}
		result = append(result, pin)
	}
	return result
}

func s2seconds(input string, def time.Duration) time.Duration {
	i := s2int(input, math.MinInt)

	if i == math.MinInt {
		return def
	}

	return time.Duration(i) * time.Second
}

func s2int(input string, def int) int {
	if len(input) == 0 {
		return def
	}
	i, err := strconv.Atoi(input)
	if err != nil {
		return def
	}
	return i
}

func s2bool(input string, def bool) bool {
	if len(input) == 0 {
		return def
	}
	i, err := strconv.ParseBool(input)
	if err != nil {
		return def
	}
	return i
}

func s2float32(input string, def float32) float32 {
	if len(input) == 0 {
		return def
	}
	i, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return def
	}
	return float32(i)
}
