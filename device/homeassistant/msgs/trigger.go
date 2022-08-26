package msgs

type StateTriggerRequest struct {
	BaseMessage
	Trigger StateTriggerDetails `json:"trigger"`
}

type StateTriggerDetails struct {
	Platform string  `json:"platform"`
	EntityID string  `json:"entity_id"`
	From     *string `json:"from,omitempty"`
	To       *string `json:"to,omitempty"`
}
