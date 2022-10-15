package config

import "strings"

const (
	cSystemName = "general.name"

	cExitDelay            = "general.exitdelay"
	cEntryDelay           = "general.entrydelay"
	cMaxWrongPinCount     = "general.maxwrongpincount"
	cBatteryThreshold     = "general.batterythreshold"
	cLinkQualityThreshold = "general.linkqualitythreshold"

	cTamperAlarmWhileArmed    = "general.tamperalarmwhilearmed"
	cTamperAlarmWhileDisarmed = "general.tamperalarmwhiledisarmed"

	cDeviceSystemFaultAlarm      = "general.devicesystemfaultalarm"
	cDeviceSystemFaultAlarmDelay = "general.devicesystemfaultalarmdelay"

	cDevices = "devices"
	cZones   = "zones"
	cPersons = "persons"

	cZigbee2MqttConfigs   = "connectors.zigbee2mqtt"
	cHomeAssistantConfigs = "connectors.homeassistant"

	cLogLevel  = "log.level"
	cLogFormat = "log.format"
)

func Init() {
	initViper()
	setDefault(cSystemName, "Home")
	setDefault(cExitDelay, 60)
	setDefault(cEntryDelay, 60)
	setDefault(cMaxWrongPinCount, 3)
	setDefault(cBatteryThreshold, 0.1)
	setDefault(cLinkQualityThreshold, 0.1)
	setDefault(cTamperAlarmWhileArmed, true)
	setDefault(cTamperAlarmWhileDisarmed, false)
	setDefault(cDeviceSystemFaultAlarm, true)
	setDefault(cDeviceSystemFaultAlarmDelay, 600)

	setDefault(cDevices, []DeviceConfig{})
	setDefault(cZones, []ZoneConfig{})
	setDefault(cPersons, []Person{})
	setDefault(cZigbee2MqttConfigs, []Zigbee2MqttConfig{})
	setDefault(cHomeAssistantConfigs, []HomeAssistantConfig{})

	setDefault(cLogLevel, "info")
	setDefault(cLogFormat, "text")
}

func SystemName() string {
	return getString(cSystemName)
}

func ExitDelay() int {
	return getInt(cExitDelay)
}

func EntryDelay() int {
	return getInt(cEntryDelay)
}

func MaxWrongPinCount() int {
	return getInt(cMaxWrongPinCount)
}

func BatteryLevelThreshold() float32 {
	return getFloat32(cBatteryThreshold)
}

func LinkQualityThreshold() float32 {
	return getFloat32(cLinkQualityThreshold)
}

func DeviceSystemFaultAlarm() bool {
	return getBool(cDeviceSystemFaultAlarm)
}

func DeviceSystemFaultDelay() int {
	return getInt(cDeviceSystemFaultAlarmDelay)
}

func TamperAlarmWhileArmed() bool {
	return getBool(cTamperAlarmWhileArmed)
}

func TamperAlarmWhileDisarmed() bool {
	return getBool(cTamperAlarmWhileDisarmed)
}

func DeviceConfigs() []DeviceConfig {
	return getObjects[DeviceConfig](cDevices)
}

func ZoneConfigs() []ZoneConfig {
	return getObjects[ZoneConfig](cZones)
}

func Persons() []Person {
	return getObjects[Person](cPersons)
}

func Zigbee2MqttConfigs() []Zigbee2MqttConfig {
	return getObjects[Zigbee2MqttConfig](cZigbee2MqttConfigs)
}

func HomeAssistantConfigs() []HomeAssistantConfig {
	return getObjects[HomeAssistantConfig](cHomeAssistantConfigs)
}

func LogFormat() string {
	return strings.ToLower(getString(cLogFormat))
}

func LogLevel() string {
	return strings.ToLower(getString(cLogLevel))
}
