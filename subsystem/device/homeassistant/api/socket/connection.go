package socket

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/msgs"
	"sync"
	"sync/atomic"
	"time"
)

type subscription struct {
	channel     chan []byte
	origPayload SetId
}

type Connection struct {
	socket        *Socket
	commands      sync.Map
	subscriptions sync.Map
	seq           uint64
	writerChan    chan interface{}
}

func NewConnection() *Connection {
	return &Connection{
		socket:        NewSocket(),
		commands:      sync.Map{},
		subscriptions: sync.Map{},
		seq:           0,
		writerChan:    nil,
	}
}

func (c *Connection) nextSeq() uint64 {
	return atomic.AddUint64(&c.seq, 1)
}

func (c *Connection) Connect(url string, token string) error {
	ch, err := c.socket.Connect(url)
	if err != nil {
		return err
	}
	if ch != nil {
		c.commands.Range(func(key, value any) bool {
			value.(chan Data) <- Data{
				Id:  0,
				Msg: nil,
				Err: remoteError{
					Code:    "connection_closed",
					Message: "connection is closed.",
				},
			}
			return true
		})
		c.commands = sync.Map{}
		c.writerChan = make(chan interface{})
		go c.readPump(url, token, ch)
		c.subscriptions.Range(func(key, value any) bool {
			c.Subscribe(value.(SetId), nil)
			return true
		})
	}
	return nil
}

func (c *Connection) readPump(url string, token string, ch chan []byte) {
	for data := range ch {
		var result msgs.BaseResult
		err := json.Unmarshal(data, &result)
		if err != nil {
			log.Error().Err(err).Msg("Could not parse json")
			continue
		}

		switch result.Type {
		case msgs.AuthRequired:
			c.socket.SendJson(msgs.AuthRequest{
				Type:        msgs.Auth,
				AccessToken: token,
			})
		case msgs.AuthInvalid:
			log.Error().Msg("Authentication is invalid")
		case msgs.AuthOk:
			log.Info().Msg("Authentication successful")
			go c.writerPump()
		default:
			dataResult := Data{
				Id:  result.Id,
				Msg: data,
				Err: nil,
			}
			if result.Success != nil && *(result.Success) == false {
				dataResult.Err = remoteError{
					Code:    result.Error.Code,
					Message: result.Error.Message,
				}
			}

			ch, ok := c.commands.LoadAndDelete(result.Id)
			if ok {
				cc := ch.(chan Data)
				cc <- dataResult
				continue
			}

			ch, ok = c.subscriptions.Load(result.Id)
			if ok {
				sub := ch.(subscription)
				sub.channel <- dataResult.Msg
				continue
			}

			log.Debug().Uint64("id", result.Id).Msg("No handler registered for message.")
		}
	}
	log.Info().Str("url", url).Msg("Connection closed. Reconnecting in 10 seconds...")
	time.AfterFunc(10*time.Second, func() {
		c.Connect(url, token)
	})
}

func (c *Connection) Disconnect() {
	close(c.writerChan)
	c.Disconnect()
}

func (c *Connection) sendJson(json interface{}) {
	c.writerChan <- json
}

func (c *Connection) writerPump() {
	if c.writerChan == nil {
		return
	}

	for data := range c.writerChan {
		c.socket.SendJson(data)
	}
}

func (c *Connection) Subscribe(payload SetId, result interface{}) (chan []byte, uint64, error) {
	seq := c.nextSeq()
	payload.SetId(seq)

	subs := subscription{
		channel:     make(chan []byte),
		origPayload: payload,
	}
	c.subscriptions.Store(seq, subs)

	err := c.cmd(seq, payload, result)
	return subs.channel, seq, err
}

func (c *Connection) Unsubscribe(id uint64) {
	ch, ok := c.subscriptions.Load(id)
	if ok {
		sub := ch.(subscription)
		close(sub.channel)
		c.subscriptions.Delete(id)
	}
}

func (c *Connection) Command(payload SetId, result interface{}) error {
	seq := c.nextSeq()
	payload.SetId(seq)
	return c.cmd(seq, payload, result)
}

func (c *Connection) cmd(id uint64, payload SetId, result interface{}) error {
	rch := make(chan Data)
	c.commands.Store(id, rch)

	c.sendJson(payload)

	data := <-rch

	if data.Err != nil {
		return data.Err
	}

	if result != nil {
		err := json.Unmarshal(data.Msg, &result)
		return err
	}

	return nil
}
