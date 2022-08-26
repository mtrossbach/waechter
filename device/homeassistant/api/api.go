package api

import (
	"github.com/mtrossbach/waechter/device/homeassistant/api/socket"
	msgs2 "github.com/mtrossbach/waechter/device/homeassistant/msgs"
)

type Api struct {
	url        string
	token      string
	connection *socket.Connection
}

func NewApi(url string, token string) *Api {
	return &Api{
		url:        url,
		token:      token,
		connection: socket.NewConnection(),
	}
}

func (a *Api) Connect() {
	a.connection.Connect(a.url, a.token)
}

func (a *Api) Disconnect() {
	a.connection.Disconnect()
}

func (a *Api) GetStates() (msgs2.StateResult, error) {
	var result msgs2.StateResult
	err := a.executeBasicCommand(msgs2.GetStates, &result)
	return result, err
}

func (a *Api) Ping() (msgs2.BaseMessage, error) {
	var result msgs2.BaseMessage
	err := a.executeBasicCommand(msgs2.Ping, &result)
	return result, err
}

func (a *Api) SubscribeStateTrigger(entityId string) (interface{}, chan []byte, uint64, error) {
	payload := msgs2.StateTriggerRequest{
		BaseMessage: msgs2.BaseMessage{
			Type: msgs2.SubscribeTrigger,
		},
		Trigger: msgs2.StateTriggerDetails{
			Platform: "state",
			EntityID: entityId,
		},
	}
	var result interface{}
	ch, id, err := a.connection.Subscribe(&payload, &result)
	return result, ch, id, err
}

func (a *Api) UnsubscribeStateTrigger(id uint64) (msgs2.BaseResult, error) {
	a.connection.Unsubscribe(id)
	payload := msgs2.UnsubscribeRequest{
		BaseMessage: msgs2.BaseMessage{
			Type: msgs2.UnsubscribeEvents,
		},
		Subscription: id,
	}
	var result msgs2.BaseResult
	err := a.connection.Command(&payload, &result)
	return result, err
}

func (a *Api) executeBasicCommand(mtype msgs2.MsgType, result interface{}) error {
	payload := msgs2.BaseMessage{
		Type: mtype,
	}
	err := a.connection.Command(&payload, &result)
	return err
}
