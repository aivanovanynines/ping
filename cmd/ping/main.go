package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aivanovanynines/ping/internal/app"
)

func main() {
	application := app.New(8080)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-stopChan
		log.Print("interrupt signal")
		cancel()
	}()

	application.Run(ctx)
}
