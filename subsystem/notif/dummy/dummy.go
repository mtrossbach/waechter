package dummy

import (
	"fmt"
	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type dummy struct {
	log zerolog.Logger
}

func New() *dummy {
	return &dummy{
		log: misc.Logger("DummyNotif"),
	}
}

func (d *dummy) GetName() string {
	return "DummyNotif"
}

func (d *dummy) SendNotification(notif system.Notification) {
	d.log.Info().Str("type", string(notif.Type)).Str("title", notif.Title).Msg(fmt.Sprintf("##### %v #####", notif.Description))
}
