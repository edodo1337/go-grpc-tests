package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := NewApp()
	defer app.Shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Server run error: %v", err)

			return
		}
	}()

	<-stop
}
