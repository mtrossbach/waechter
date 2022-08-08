package boot

import (
	"github.com/mtrossbach/waechter/waechter/subsystem/zigbee2mqtt"
	"github.com/mtrossbach/waechter/waechter/system"
)

func Boot() {
	system := system.NewWaechterSystem()
	system.RegisterSubsystem(zigbee2mqtt.NewZigbee2MqttSubsystem())
}
