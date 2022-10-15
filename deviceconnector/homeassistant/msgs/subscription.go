package msgs

type SubscribeStateRequest struct {
	BaseMessage
	Trigger TriggerDetails `json:"trigger"`
}

type TriggerDetails struct {
	Platform string  `json:"platform"`
	EntityID string  `json:"entity_id"`
	From     *string `json:"from,omitempty"`
	To       *string `json:"to,omitempty"`
}

type UnsubscribeRequest struct {
	BaseMessage
	Subscription uint64 `json:"subscription"`
}
