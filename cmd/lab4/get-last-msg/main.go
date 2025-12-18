package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"ris/internal/subscriber"

	"github.com/nats-io/nats.go"
)

func main() {
	msgType := flag.String("type", "prize", "Message type: 'prize' or 'laureate'")
	flag.Parse()

	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer natsConn.Close()

	subs := subscriber.New(natsConn)
	defer subs.Close()

	switch *msgType {
	case "prize":
		prize, err := subs.GetLastPrizeMessage()
		if err != nil {
			slog.Error("Failed to get last prize message", "error", err)
			os.Exit(1)
		}
		fmt.Println("Last Prize Message:")
		fmt.Printf("Year: %s\n", prize.Year)
		fmt.Printf("Category: %s\n", prize.Category)
		fmt.Printf("Overall Motivation: %s\n", prize.OverallMotivation)
		fmt.Printf("Number of Laureates: %d\n", len(prize.Laureates))
		for i, laureate := range prize.Laureates {
			fmt.Printf("  Laureate %d: %s %s\n", i+1, laureate.Firstname, laureate.Surname)
		}

	case "laureate":
		laureate, err := subs.GetLastLaureateMessage()
		if err != nil {
			slog.Error("Failed to get last laureate message", "error", err)
			os.Exit(1)
		}
		fmt.Println("Last Laureate Message:")
		fmt.Printf("ID: %d\n", laureate.Id)
		fmt.Printf("Name: %s %s\n", laureate.Firstname, laureate.Surname)
		fmt.Printf("Motivation: %s\n", laureate.Motivation)
		fmt.Printf("Share: %d\n", laureate.Share)

	default:
		slog.Error("Invalid message type", "type", *msgType)
		os.Exit(1)
	}

	slog.Info("Successfully retrieved last message")
}
