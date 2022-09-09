package driver

import (
	"encoding/json"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/homeassistant/connector"
	"github.com/mtrossbach/waechter/device/homeassistant/msgs"
	"github.com/mtrossbach/waechter/system"
)

func ContactSensorHandler(dev *system.Device, controller device.SystemController) connector.MessageHandler {
	return func(msgType msgs.MsgType, msg []byte) {
		if msgType != msgs.Event {
			return
		}

		var event msgs.EventResponse
		err := json.Unmarshal(msg, &event)
		if err != nil {
			return
		}
		if event.Event.Variables.Trigger.ToState.State == "on" {
			controller.Alarm(system.BurglarAlarm, *dev)
		}
	}
}
