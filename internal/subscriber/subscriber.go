package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ris/internal/domain"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	subjectPrizeCreated    = "prize.created"
	subjectLaureateCreated = "laureate.created"
)

type Subscriber struct {
	conn *nats.Conn
	js   nats.JetStreamContext

	subs []*nats.Subscription
}

func New(conn *nats.Conn) *Subscriber {
	js, _ := conn.JetStream()
	_, _ = js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{subjectPrizeCreated, subjectLaureateCreated},
		MaxAge:   time.Hour * 24 * 7, // 7 days
	})

	return &Subscriber{
		conn: conn,
		js:   js,
		subs: make([]*nats.Subscription, 0),
	}
}

func (s *Subscriber) Close() error {
	errs := make([]error, 0, len(s.subs))
	for _, sub := range s.subs {
		errs = append(errs, sub.Unsubscribe())
	}
	return errors.Join(errs...)
}

func (s *Subscriber) SubscribePrizeCreated(handler func(prize domain.Prize) error) error {
	sub, err := s.conn.Subscribe(subjectPrizeCreated, func(msg *nats.Msg) {
		var prize domain.Prize
		if err := json.Unmarshal(msg.Data, &prize); err != nil {
			slog.Error("failed to unmarshal prize created event: %w", err)
			return
		}
		if err := handler(prize); err != nil {
			slog.Error("failed to handle prize created event: %w", err)
			msg.Nak()
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to prize created event: %w", err)
	}
	s.subs = append(s.subs, sub)
	return nil
}

func (s *Subscriber) SubscribeLaureateCreated(handler func(laureate domain.Laureate) error) error {
	sub, err := s.conn.Subscribe(subjectLaureateCreated, func(msg *nats.Msg) {
		var laureate domain.Laureate
		if err := json.Unmarshal(msg.Data, &laureate); err != nil {
			slog.Error("failed to unmarshal laureate created event: %w", err)
			return
		}
		if err := handler(laureate); err != nil {
			slog.Error("failed to handle laureate created event: %w", err)
			msg.Nak()
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to laureate created event: %w", err)
	}
	s.subs = append(s.subs, sub)
	return nil
}

func (s *Subscriber) GetLastPrizeMessage() (*domain.Prize, error) {
	msg, err := s.js.GetLastMsg("EVENTS", subjectPrizeCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get last prize message: %w", err)
	}
	var prize domain.Prize
	if err := json.Unmarshal(msg.Data, &prize); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prize: %w", err)
	}
	return &prize, nil
}

func (s *Subscriber) GetLastLaureateMessage() (*domain.Laureate, error) {
	msg, err := s.js.GetLastMsg("EVENTS", subjectLaureateCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to get last laureate message: %w", err)
	}
	var laureate domain.Laureate
	if err := json.Unmarshal(msg.Data, &laureate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal laureate: %w", err)
	}
	return &laureate, nil
}
