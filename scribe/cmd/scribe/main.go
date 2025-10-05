package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	daemon := flag.Bool("daemon", true, "run as background agent")
	flag.Parse()

	log.Printf("scribe starting (daemon=%v)", *daemon)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Minimal work loop placeholder
	loop := time.NewTicker(2 * time.Second)
	defer loop.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("scribe stopped")
			return
		case <-loop.C:
			log.Printf("scribe heartbeat")
		}
	}
}
