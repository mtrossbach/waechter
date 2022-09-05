package whatsapp

type MessagePayload struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type Language struct {
	Code string `json:"code"`
}
type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}
type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}
