package runner

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mtrossbach/waechter/misc"
	"github.com/mtrossbach/waechter/subsystem/zigbee2mqtt"
	"github.com/mtrossbach/waechter/system"
)

func Run() {
	misc.InitializeLogging()
	misc.Log.Info("Starting up...")
	system := system.NewWaechterSystem()
	system.RegisterSubsystem(zigbee2mqtt.NewZigbee2MqttSubsystem())
	misc.Log.Info("Started.")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Printf("Caught SIGTERM %v", sig)
}
