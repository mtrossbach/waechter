package model

type ZigbeeDevice struct {
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
