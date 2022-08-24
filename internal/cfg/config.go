package cfg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.waechter")
	viper.AddConfigPath("/etc/waechter/")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("%w", err)
	}
	updateLogger()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		updateLogger()
		Print()
	})
	viper.WatchConfig()
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetFloat32(key string) float32 {
	return float32(viper.GetFloat64(key))
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStrings(key string) []string {
	return viper.GetStringSlice(key)
}

func SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func Print() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	fmt.Printf("########################################\n")
	for _, k := range keys {
		if strings.Contains(strings.ToLower(k), "pwd") || strings.Contains(strings.ToLower(k), "password") || strings.Contains(strings.ToLower(k), "pins") || strings.Contains(strings.ToLower(k), "token") {
			fmt.Printf("  %v: %v\n", k, "***")
		} else {
			fmt.Printf("  %v: %v\n", k, viper.Get(k))
		}
	}
	fmt.Printf("########################################\n")
}