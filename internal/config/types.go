package config

type DeviceConfig struct {
	Id   string `json:"id"`
	Zone string `json:"zone"`
}

type Zigbee2MqttConfig struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	ClientId  string `json:"clientId"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	BaseTopic string `json:"baseTopic"`
}

type HomeAssistantConfig struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

type ZoneConfig struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Perimeter   bool   `json:"perimeter"`
	Delayed     bool   `json:"delayed"`
}

type Person struct {
	Name     string `json:"name"`
	Pin      string `json:"pin"`
	Lang     string `json:"lang"`
	WhatsApp string `json:"whatsapp"`
}
