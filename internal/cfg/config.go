package cfg

import (
	"fmt"
	"sort"
	"strings"

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

func GetStringStringMaps(key string) []map[string]string {
	data := viper.Get(key).([]interface{})
	var result []map[string]string

	for _, d := range data {
		dm := d.(map[string]interface{})

		m := make(map[string]string)

		for k, v := range dm {
			m[k] = fmt.Sprintf("%v", v)
		}
		result = append(result, m)
	}

	return result
}

func SetDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func SetString(key string, value string) {
	viper.Set(key, value)
}

func WriteConfig() {
	viper.WriteConfig()
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
