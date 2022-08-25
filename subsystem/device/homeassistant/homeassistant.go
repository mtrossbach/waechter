package homeassistant

import (
	"github.com/gorilla/websocket"
	"github.com/mtrossbach/waechter/internal/cfg"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/api"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/device"
	"github.com/mtrossbach/waechter/subsystem/device/homeassistant/model"
	"github.com/mtrossbach/waechter/system"
)

type homeassistant struct {
	c   *websocket.Conn
	api *api.Api
}

func New() *homeassistant {
	ins := &homeassistant{
		api: api.NewApi(cfg.GetString("homeassistant.url"), cfg.GetString("homeassistant.token")),
	}
	return ins
}

func (ha *homeassistant) GetName() string {
	return model.SubsystemName
}

func (ha *homeassistant) Start(system system.DeviceSystem) {
	ha.api.Connect()

	st, _ := ha.api.GetStates()
	for _, s := range st.Result {
		if s.Attributes.DeviceClass == "motion" && s.Attributes.MotionValid {
			dev := device.NewMotionSensor(ha.api, s.EntityID, s.Attributes.FriendlyName)
			system.AddDevice(dev)
		}
	}
}

func (ha *homeassistant) Stop() {

}
