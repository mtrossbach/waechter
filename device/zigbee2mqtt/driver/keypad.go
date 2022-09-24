package driver

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
)

func KeypadHandler(dev *system.Device, controller device.SystemController, sender Sender) connector.MessageHandler {
	return func(msg mqtt.Message) {
		var payload keypadStatus
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
			return
		}

		log.Debug().RawJSON("payload", msg.Payload()).Msg("Got data")

		if payload.Battery > 0 {
			level := float32(payload.Battery) / float32(100)
			controller.ReportBatteryLevel(level, *dev)
		}

		if payload.LinkQuality > 0 {
			controller.ReportLinkQuality(float32(payload.LinkQuality)/float32(255), *dev)
		}

		if payload.Tamper {
			controller.Alarm(system.TamperAlarm, *dev)
		}

		if len(payload.Action) > 0 {
			keypadSendState(payload.Action, &payload.ActionTransaction, sender) //Send confirmation (required for some devices)

			if payload.Action == "arm_day_zones" {
				controller.ArmStay(*dev)
			} else if payload.Action == "arm_all_zones" {
				controller.ArmAway(*dev)
			} else if payload.Action == "disarm" {
				controller.Disarm(payload.ActionCode, *dev)
			} else if payload.Action == "panic" {
				controller.Alarm(system.PanicAlarm, *dev)
			}
		}
	}
}

func keypadSendState(state string, transactionId *int, sender Sender) {
	sender(newKeypadSetState(state, transactionId))
}

func KeypadStateUpdater(controller device.SystemController, sender Sender) {
	keypadSendState(sysStateToKeypadState(controller.GetArmState(), controller.GetAlarmType()), nil, sender)
}
