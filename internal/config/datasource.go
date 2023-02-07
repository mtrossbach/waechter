package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"path"
	"sort"
	"strings"
)

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.waechter")
	viper.AddConfigPath("/etc/waechter/")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("%w", err)
	}
}

func getBool(key string) bool {
	return viper.GetBool(key)
}

func getFloat32(key string) float32 {
	return float32(viper.GetFloat64(key))
}

func getInt(key string) int {
	return viper.GetInt(key)
}

func getString(key string) string {
	return viper.GetString(key)
}

func getStrings(key string) []string {
	return viper.GetStringSlice(key)
}

func getObject[T any](key string) *T {
	data := viper.Get(key)
	if data == nil {
		return nil
	}
	data = data.(interface{})

	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var r T
	json.Unmarshal(b, &r)
	if err != nil {
		return nil
	}

	return &r
}

func getObjects[T any](key string) []T {
	data := viper.Get(key)
	if data == nil {
		return []T{}
	}
	data = data.(interface{})

	b, err := json.Marshal(data)
	if err != nil {
		return []T{}
	}

	r := make([]T, 0)
	json.Unmarshal(b, &r)
	if err != nil {
		return []T{}
	}

	return r
}

func getStringStringMaps(key string) []map[string]string {
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

func setDefault(key string, value interface{}) {
	viper.SetDefault(key, value)
}

func ConfigFileDir() string {
	return path.Dir(viper.ConfigFileUsed())
}

func Print() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	fmt.Printf("########################################\n")
	for _, k := range keys {
		if strings.Contains(strings.ToLower(k), "pwd") || strings.Contains(strings.ToLower(k), "password") || strings.Contains(strings.ToLower(k), "pin") || strings.Contains(strings.ToLower(k), "token") {
			fmt.Printf("  %v: %v\n", k, "***")
		} else {
			fmt.Printf("  %v: %v\n", k, viper.Get(k))
		}
	}
	fmt.Printf("########################################\n")
}
