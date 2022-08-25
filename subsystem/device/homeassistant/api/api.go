package api

import (
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/api/socket"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/msgs"
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

func (a *Api) GetStates() (msgs.StateResult, error) {
	var result msgs.StateResult
	err := a.executeBasicCommand(msgs.GetStates, &result)
	return result, err
}

func (a *Api) Ping() (msgs.BaseMessage, error) {
	var result msgs.BaseMessage
	err := a.executeBasicCommand(msgs.Ping, &result)
	return result, err
}

func (a *Api) SubscribeStateTrigger(entityId string) (interface{}, chan []byte, uint64, error) {
	payload := msgs.StateTriggerRequest{
		BaseMessage: msgs.BaseMessage{
			Type: msgs.SubscribeTrigger,
		},
		Trigger: msgs.StateTriggerDetails{
			Platform: "state",
			EntityID: entityId,
		},
	}
	var result interface{}
	ch, id, err := a.connection.Subscribe(&payload, &result)
	return result, ch, id, err
}

func (a *Api) UnsubscribeStateTrigger(id uint64) (msgs.BaseResult, error) {
	a.connection.Unsubscribe(id)
	payload := msgs.UnsubscribeRequest{
		BaseMessage: msgs.BaseMessage{
			Type: msgs.UnsubscribeEvents,
		},
		Subscription: id,
	}
	var result msgs.BaseResult
	err := a.connection.Command(&payload, &result)
	return result, err
}

func (a *Api) executeBasicCommand(mtype msgs.MsgType, result interface{}) error {
	payload := msgs.BaseMessage{
		Type: mtype,
	}
	err := a.connection.Command(&payload, &result)
	return err
}
