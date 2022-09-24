package driver

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mtrossbach/waechter/device"
	"github.com/mtrossbach/waechter/device/zigbee2mqtt/connector"
	"github.com/mtrossbach/waechter/internal/log"
	"github.com/mtrossbach/waechter/system"
	"time"
)

func SirenHandler(dev *system.Device, controller device.SystemController) connector.MessageHandler {
	return func(msg mqtt.Message) {
		var payload sirenStatus
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			log.Error().Str("payload", string(msg.Payload())).Msg("Could not parse payload")
			return
		}

		log.Debug().Str("payload", string(msg.Payload())).Msg("Got data")

		if payload.Battery > 0 {
			controller.ReportBatteryLevel(float32(payload.Battery)/float32(100), *dev)
		}

		if payload.LinkQuality > 0 {
			controller.ReportLinkQuality(float32(payload.LinkQuality)/float32(255), *dev)
		}

		if payload.Tamper {
			controller.Alarm(system.TamperAlarm, *dev)
		}
	}
}

func SirenStateUpdater(controller device.SystemController, sender Sender) {
	sender(newWarningPayload(controller.GetAlarmType()))
	time.AfterFunc(100*time.Millisecond, func() {
		// Resend after 100ms because some sirens do not correctly process payloads during active alarms
		sender(newWarningPayload(controller.GetAlarmType()))
	})
}
