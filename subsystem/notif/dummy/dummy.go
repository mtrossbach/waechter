package dummy

import (
	"fmt"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
)

type dummy struct {
}

func New() *dummy {
	return &dummy{}
}

func (d *dummy) GetName() string {
	return "DummyNotif"
}

func (d *dummy) SendNotification(notif system.Notification) {
	log.Info().Str("type", string(notif.Type)).Str("title", notif.Title).Msg(fmt.Sprintf("##### %v #####", notif.Description))
}
