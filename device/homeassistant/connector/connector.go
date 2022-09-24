package connector

import (
	"encoding/json"
	"fmt"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/internal/log"
	"golang.org/x/net/websocket"
	"sync"
	"sync/atomic"
	"time"
)

type DisconnectedHandler func(err error)
type MessageHandler func(msgType msgs.MsgType, msg []byte)

type Connector struct {
	ws         *websocket.Conn
	cmd        sync.Map
	sub        sync.Map
	handler    DisconnectedHandler
	connection bool
	seq        uint64
	writerChan chan any
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) Connect(url string, token string, handler DisconnectedHandler) error {
	c.handler = handler
	c.sub = sync.Map{}
	c.cmd = sync.Map{}
	c.seq = 0
	c.writerChan = make(chan any)

	ws, err := websocket.Dial(url, "", url)
	if err != nil {
		return err
	}
	c.ws = ws
	c.connection = true
	go c.readerPump(token)

	return nil
}

func (c *Connector) Disconnect() {
	c.disconnect(nil)
}

func (c *Connector) disconnect(err error) {
	if !c.connection {
		return
	}
	c.connection = false
	if c.writerChan != nil {
		close(c.writerChan)
		c.writerChan = nil
	}
	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}
	if c.handler != nil {
		c.handler(err)
		c.handler = nil
	}
}

func (c *Connector) nextSeq() uint64 {
	return atomic.AddUint64(&c.seq, 1)
}

func (c *Connector) readerPump(token string) {
	for c.connection {
		var payload []byte
		err := websocket.Message.Receive(c.ws, &payload)
		if err != nil {
			c.disconnect(err)
			return
		}
		c.handleMessage(payload, token)
	}
}

func (c *Connector) handleMessage(msg []byte, token string) {
	var result msgs.BaseResult
	err := json.Unmarshal(msg, &result)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse json")
		return
	}

	switch result.Type {
	case msgs.AuthRequired:
		c.send(msgs.AuthRequest{
			Type:        msgs.Auth,
			AccessToken: token,
		})
	case msgs.AuthInvalid:
		log.Error().Msg("HomeAssistant authentication is invalid")
	case msgs.AuthOk:
		log.Debug().Msg("HomeAssistant authentication successful")
		go c.writerPump()
	default:
		ch, ok := c.cmd.Load(result.Id)
		if ok {
			cc := ch.(chan []byte)
			cc <- msg
			return
		}

		ch, ok = c.sub.Load(result.Id)
		if ok {
			sub := ch.(MessageHandler)
			sub(result.Type, msg)
			return
		}

		log.Debug().Uint64("id", result.Id).Msg("No handler registered for message.")
	}
}

func (c *Connector) command(id uint64, payload any, result any) error {
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
		c.disconnect(err)
		return err
	}
}

func (c *Connector) basicCommand(id uint64, payload SetId) error {
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

func (c *Connector) Command(payload SetId, result any) error {
	id := c.nextSeq()
	payload.SetId(id)
	return c.command(id, payload, result)
}

func (c *Connector) Subscribe(payload SetId, handler MessageHandler) (uint64, error) {
	id := c.nextSeq()
	payload.SetId(id)
	c.sub.Store(id, handler)

	err := c.basicCommand(id, payload)
	return id, err
}

func (c *Connector) Unsubscribe(id uint64) {
	c.sub.Delete(id)
}

func (c *Connector) SubscribeStateTrigger(entityId string, handler MessageHandler) (uint64, error) {
	payload := msgs.SubscribeStateRequest{
		BaseMessage: msgs.BaseMessage{
			Type: msgs.SubscribeTrigger,
		},
		Trigger: msgs.TriggerDetails{
			Platform: "state",
			EntityID: entityId,
		},
	}
	id, err := c.Subscribe(&payload, handler)
	return id, err
}

func (c *Connector) UnsubscribeStateTrigger(id uint64) error {
	c.Unsubscribe(id)
	payload := msgs.UnsubscribeRequest{
		BaseMessage: msgs.BaseMessage{
			Type: msgs.UnsubscribeEvents,
		},
		Subscription: id,
	}
	seqId := c.nextSeq()
	payload.SetId(seqId)
	return c.basicCommand(seqId, &payload)
}

func (c *Connector) writerPump() {
	if c.writerChan == nil {
		return
	}

	for data := range c.writerChan {
		c.send(data)
	}
}

func (c *Connector) send(payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = websocket.Message.Send(c.ws, string(data))
	return err
}
