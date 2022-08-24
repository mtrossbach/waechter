package homeassistant

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/api"
	"github.com/mtrossbach/waechter/system"
	"github.com/rs/zerolog"
)

type homeassistant struct {
	log zerolog.Logger
	c   *websocket.Conn
	api *api.Api
}

func New() *homeassistant {
	ins := &homeassistant{
		log: cfg.Logger("HomeAssistant"),
		api: api.NewApi(cfg.GetString("homeassistant.url"), cfg.GetString("homeassistant.token")),
	}
	return ins
}

func (ha *homeassistant) GetName() string {
	return "HomeAssistant"
}

func (ha *homeassistant) Start(system system.DeviceSystem) {
	ha.api.Connect()

	go ha.test()
}

func (ha *homeassistant) test() {
	r, c, _, _ := ha.api.SubscribeStateTrigger("binary_sensor.badezimmer_sensor_motion")

	fmt.Printf("%v\n", r)
	for e := range c {
		fmt.Printf("%v\n", e)
	}
}

func (ha *homeassistant) Stop() {

}
