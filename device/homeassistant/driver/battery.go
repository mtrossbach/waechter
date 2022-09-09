package driver

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/homeassistant/connector"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/system"
	"strconv"
)

func BatteryHandler(dev *system.Device, controller device.SystemController) connector.MessageHandler {
	return func(msgType msgs.MsgType, msg []byte) {
		if msgType != msgs.Event {
			return
		}

		var event msgs.EventResponse
		err := json.Unmarshal(msg, &event)
		if err != nil {
			return
		}

		level, err := strconv.Atoi(event.Event.Variables.Trigger.ToState.State)
		if err != nil {
			system.DError(dev).Err(err).Msg("Could not parse battery level")
		} else {
			controller.ReportBatteryLevel(float32(level)/100.0, *dev)
		}
	}
}
