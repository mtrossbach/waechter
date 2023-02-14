package zigbee2mqtt

type Z2MDeviceInfo struct {
	IeeeAddress        string     `json:"ieee_address"`
	Type               string     `json:"type"`
	Supported          bool       `json:"supported"`
	FriendlyName       string     `json:"friendly_name"`
	Definition         Definition `json:"definition"`
	PowerSource        string     `json:"power_source"`
	DateCode           string     `json:"date_code"`
	ModelId            string     `json:"model_id"`
	Interviewing       bool       `json:"interviewing"`
	InterviewCompleted bool       `json:"interview_completed"`
	Manufacturer       string     `json:"manufacturer"`
}

type Clusters struct {
	Input  []string      `json:"input"`
	Output []interface{} `json:"output"`
}
type Definition struct {
	Model       string    `json:"model"`
	Vendor      string    `json:"vendor"`
	Description string    `json:"description"`
	Options     []Options `json:"options"`
	Exposes     []Exposes `json:"exposes"`
}

type Exposes struct {
	Name     string `json:"name"`
	Property string `json:"property"`
	Type     string `json:"type"`
}
type Options struct {
	Name     string `json:"name"`
	Property string `json:"property"`
	Type     string `json:"type"`
}

type DeviceEvent struct {
	Data Data   `json:"data"`
	Type string `json:"type"`
}
type Data struct {
	FriendlyName string `json:"friendly_name"`
	IeeeAddress  string `json:"ieee_address"`
}
