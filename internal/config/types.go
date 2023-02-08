package config

type Config struct {
	General       GeneralConfig          `yaml:"general"`
	Log           LogConfig              `yaml:"log"`
	Persons       []Person               `yaml:"persons"`
	Devices       []DeviceConfig         `yaml:"devices"`
	Zones         []ZoneConfig           `yaml:"zones"`
	Zigbee2Mqtt   []Zigbee2MqttConfig    `yaml:"zigbee2mqtt"`
	HomeAssistant []HomeAssistantConfig  `yaml:"homeassistant"`
	WhatsApp      *WhatsAppConfiguration `yaml:"whatsapp"`
	Notification  []string               `yaml:"notifications"`
}

type GeneralConfig struct {
	Name                        string  `yaml:"name" default:"My Home"`
	ExitDelay                   int     `yaml:"exitDelay" default:"30"`
	EntryDelay                  int     `yaml:"entryDelay" default:"30"`
	MaxWrongPinCount            int     `yaml:"maxWrongPinCount" default:"2"`
	BatteryThreshold            float32 `yaml:"batteryThreshold" default:"0.1"`
	LinkQualityThreshold        float32 `yaml:"linkQualityThreshold" default:"0.1"`
	TamperAlarmWhileArmed       bool    `yaml:"tamperAlarmWhileArmed" default:"true"`
	TamperAlarmWhileDisarmed    bool    `yaml:"tamperAlarmWhileDisarmed" default:"false"`
	DeviceSystemFaultAlarm      bool    `yaml:"deviceSystemFaultAlarm" default:"true"`
	DeviceSystemFaultAlarmDelay int     `yaml:"deviceSystemFaultAlarmDelay" default:"300"`
}

type LogConfig struct {
	Level  string `yaml:"level" default:"info"`
	Format string `yaml:"format" default:"text"`
}

type DeviceConfig struct {
	Id   string `yaml:"id"`
	Zone string `yaml:"zone"`
}

type Zigbee2MqttConfig struct {
	Id        string `yaml:"id"`
	Url       string `yaml:"url"`
	ClientId  string `yaml:"clientId"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	BaseTopic string `yaml:"baseTopic"`
}

type HomeAssistantConfig struct {
	Id    string `yaml:"id"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type ZoneConfig struct {
	Id          string `yaml:"id"`
	DisplayName string `yaml:"displayName"`
	Perimeter   bool   `yaml:"perimeter"`
	Delayed     bool   `yaml:"delayed"`
}

type Person struct {
	Name     string `yaml:"name"`
	Pin      string `yaml:"pin"`
	Lang     string `yaml:"lang"`
	WhatsApp string `yaml:"whatsapp"`
}

type WhatsAppConfiguration struct {
	PhoneId              string `yaml:"phoneId"`
	TemplateAlarm        string `yaml:"templateAlarm"`
	TemplateAutoArm      string `yaml:"templateAutoArm"`
	TemplateAutoDisarm   string `yaml:"templateAutoDisarm"`
	TemplateNotification string `yaml:"templateNotification"`
	TemplateRecover      string `yaml:"templateRecover"`
	Token                string `yaml:"token"`
}
