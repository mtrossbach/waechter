package msgs

type UnsubscribeRequest struct {
	BaseMessage
	Subscription uint64 `json:"subscription"`
}
