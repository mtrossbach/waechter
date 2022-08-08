package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mtrossbach/waechter/waechter/boot"
	"github.com/mtrossbach/waechter/waechter/misc"
)

func main() {
	misc.InitializeLogging()
	misc.Log.Info("Starting up...")
	boot.Boot()
	misc.Log.Info("Started.")
	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	/*go func() {
		// start your software here. Maybe your need to replace the for loop with other code
		for {
			// replace the time.Sleep with your code
			log.Println("Loop tick")
			time.Sleep(time.Second)
		}
	}()*/
	sig := <-cancelChan
	log.Printf("Caught SIGTERM %v", sig)
	// shutdown other goroutines gracefully
	// close other resources
}
