package socket

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

type Socket struct {
	log        zerolog.Logger
	ws         *websocket.Conn
	connection bool
}

func NewSocket() *Socket {
	return &Socket{
		log:        cfg.Logger("HASocket"),
		ws:         nil,
		connection: false,
	}
}

func (s *Socket) Connect(url string) (chan []byte, error) {
	if !s.connection {
		s.connection = true

		s.log.Debug().Str("url", url).Msg("Connecting to server...")
		ws, err := websocket.Dial(url, "", url)
		if err != nil {
			return nil, err
		}
		s.log.Info().Str("url", url).Msg("Connected to server")
		s.ws = ws

		c := make(chan []byte)
		go s.readerPump(ws, c)
		return c, nil
	}
	return nil, nil
}

func (s *Socket) readerPump(ws *websocket.Conn, c chan []byte) {
	defer ws.Close()

	for s.connection {
		var payload []byte
		err := websocket.Message.Receive(ws, &payload)
		if err != nil {
			s.log.Error().Err(err).Msg("Error while reading from socket")
			close(c)
			return
		}
		s.log.Debug().RawJSON("payload", payload).Msg("Receive")
		c <- payload
	}
}

func (s *Socket) Disconnect() {
	if s.connection {
		s.connection = false
		s.log.Debug().Msg("Disconnected")
	}
}

func (s *Socket) SendJson(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	s.SendRaw(data)
	return nil
}

func (s *Socket) SendRaw(payload []byte) error {
	s.log.Debug().RawJSON("payload", payload).Msg("Send")
	err := websocket.Message.Send(s.ws, string(payload))
	return err
}
