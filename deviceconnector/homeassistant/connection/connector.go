package connection

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/deviceconnector/homeassistant/msgs"
	"github.com/mtrossbach/waechter/internal/log"
	"golang.org/x/net/websocket"
	"sync"
	"sync/atomic"
	"time"
)

type LostHandler func(conn *Connection, err error)
type ConnectedHandler func(conn *Connection)

type StateEventHandler func(entityId string, event msgs.EventResponse)

type Connection struct {
	ws         *websocket.Conn
	cmd        sync.Map
	evt        sync.Map
	connection bool
	seq        uint64
	writerChan chan any

	conId            uint64
	url              string
	token            string
	OnConnectionLost LostHandler
	OnConnect        ConnectedHandler
}

func NewConnection(url string, token string) *Connection {
	return &Connection{
		url:   url,
		token: token,
		conId: 0,
		evt:   sync.Map{},
	}
}

func (c *Connection) Connect() {
	c.cmd = sync.Map{}
	c.seq = 0
	c.conId += 1
	c.writerChan = make(chan any)

	ws, err := websocket.Dial(c.url, "", c.url)
	if err != nil {
		if c.OnConnectionLost != nil {
			c.OnConnectionLost(c, err)
		}
		return
	}
	c.ws = ws
	c.connection = true
	if c.OnConnect != nil {
		c.OnConnect(c)
	}
	go c.readerPump(c.token)
	go c.testConnection(c.conId)
}

func (c *Connection) Disconnect() {
	if !c.connection {
		return
	}
	c.connection = false
	if c.writerChan != nil {
		close(c.writerChan)
		c.writerChan = nil
	}
	if c.ws != nil {
		_ = c.ws.Close()
		c.ws = nil
	}
}

func (c *Connection) onConnectionLost(err error) {
	c.Disconnect()

	if c.OnConnectionLost != nil {
		c.OnConnectionLost(c, err)
	}

	<-time.After(1 * time.Second)
	c.Connect()
}

func (c *Connection) nextSeq() uint64 {
	return atomic.AddUint64(&c.seq, 1)
}

func (c *Connection) readerPump(token string) {
	for c.connection {
		var payload []byte
		err := websocket.Message.Receive(c.ws, &payload)
		if err != nil {
			c.onConnectionLost(err)
			return
		}
		c.handleMessage(payload, token)
	}
}

func (c *Connection) handleMessage(msg []byte, token string) {
	var result msgs.BaseResult
	err := json.Unmarshal(msg, &result)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse json")
		return
	}

	switch result.Type {
	case msgs.AuthRequired:
		_ = c.send(msgs.AuthRequest{
			Type:        msgs.Auth,
			AccessToken: token,
		})
	case msgs.AuthInvalid:
		log.Error().Msg("HomeAssistant authentication is invalid")
	case msgs.AuthOk:
		log.Debug().Msg("HomeAssistant authentication successful")
		go c.writerPump()
	case msgs.Event:

		var event msgs.EventResponse
		err := json.Unmarshal(msg, &event)
		if err != nil {
			log.Error().RawJSON("msg", msg).Err(err).Msg("Could not parse EventResponse")
			return
		}

		h, ok := c.evt.Load(event.Event.Variables.Trigger.EntityID)
		if ok {
			s := h.(subscription)
			s.handler(event.Event.Variables.Trigger.EntityID, event)
			return
		} else {
			log.Error().RawJSON("msg", msg).Msg("Received EventResponse but did not find a handler")
		}

	default:
		ch, ok := c.cmd.Load(result.Id)
		if ok {
			cc := ch.(chan []byte)
			cc <- msg
			return
		}

		log.Debug().Uint64("id", result.Id).Msg("No handler registered for message.")
	}
}

func (c *Connection) command(id uint64, payload any, result any) error {
	ch := make(chan []byte)
	defer close(ch)

	c.cmd.Store(id, ch)

	c.writerChan <- payload

	select {
	case data := <-ch:
		c.cmd.Delete(id)

		err := json.Unmarshal(data, &result)
		return err

	case <-time.After(30 * time.Second):
		err := fmt.Errorf("response timeout")
		c.onConnectionLost(err)
		return err
	}
}

func (c *Connection) basicCommand(id uint64, payload SetId) error {
	var result msgs.BaseResult
	err := c.command(id, payload, &result)
	if err != nil {
		return err
	}
	if result.Success != nil && *result.Success == true {
		return nil
	} else {
		return remoteError{
			Code:    result.Error.Code,
			Message: result.Error.Message,
		}
	}
}

func (c *Connection) Command(payload SetId, result any) error {
	id := c.nextSeq()
	payload.SetId(id)
	return c.command(id, payload, result)
}

func (c *Connection) SubscribeStateEvents(entityId string, handler StateEventHandler) error {
	if _, exists := c.evt.Load(entityId); exists {
		if err := c.UnsubscribeStateEvents(entityId); err != nil {
			return err
		}
	}

	seqId := c.nextSeq()
	payload := msgs.SubscribeStateRequest{
		BaseMessage: msgs.BaseMessage{
			Type: msgs.SubscribeTrigger,
			Id:   seqId,
		},
		Trigger: msgs.TriggerDetails{
			Platform: "state",
			EntityID: entityId,
		},
	}
	c.evt.Store(entityId, subscription{
		handler: handler,
		seqId:   seqId,
	})

	err := c.basicCommand(seqId, &payload)
	return err
}

func (c *Connection) UnsubscribeStateEvents(entityId string) error {
	if s, exists := c.evt.Load(entityId); exists {
		payload := msgs.UnsubscribeRequest{
			BaseMessage: msgs.BaseMessage{
				Type: msgs.UnsubscribeEvents,
			},
			Subscription: s.(subscription).seqId,
		}

		c.evt.Delete(entityId)
		seqId := c.nextSeq()
		payload.SetId(seqId)
		return c.basicCommand(seqId, &payload)
	}
	return nil
}

func (c *Connection) writerPump() {
	if c.writerChan == nil {
		return
	}

	for data := range c.writerChan {
		if err := c.send(data); err != nil {
			log.Error().Err(err).Interface("data", data).Msg("Could not send data to HomeAssistant")
		}
	}
}

func (c *Connection) send(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = websocket.Message.Send(c.ws, string(data))
	return err
}

func (c *Connection) testConnection(conId uint64) {
	for c.conId == conId {
		_ = c.Command(&msgs.BaseMessage{Type: msgs.Ping}, nil)
		<-time.After(30 * time.Second)
	}
}
