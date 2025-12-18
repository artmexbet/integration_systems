package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"ris/internal/domain"
	"ris/internal/subscriber"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer natsConn.Close()
	subs := subscriber.New(natsConn)
	defer subs.Close()

	err = subs.SubscribePrizeCreated(func(prize domain.Prize) error {
		prizeJSON, _ := json.MarshalIndent(prize, "", "  ")
		slog.Info("Received prize created event", "prize", string(prizeJSON))
		return nil
	})
	if err != nil {
		slog.Error("Failed to subscribe to prize created events", "error", err)
		return
	}

	err = subs.SubscribeLaureateCreated(func(laureate domain.Laureate) error {
		laureateJSON, _ := json.MarshalIndent(laureate, "", "  ")
		slog.Info("Received laureate created event", "laureate", string(laureateJSON))
		return nil
	})
	if err != nil {
		slog.Error("Failed to subscribe to laureate created events", "error", err)
		return
	}

	slog.Info("Listening for events from stream...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Server stopped")
}
