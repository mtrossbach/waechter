package msgs

type BaseMessage struct {
	Id   uint64  `json:"id"`
	Type MsgType `json:"type"`
}

func (r *BaseMessage) SetId(seq uint64) {
	r.Id = seq
}

type BaseResult struct {
	BaseMessage
	Success *bool      `json:"success"`
	Error   *BaseError `json:"error"`
}

type BaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
