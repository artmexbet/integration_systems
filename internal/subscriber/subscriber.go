package subscriber

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"ris/internal/domain"

	"github.com/nats-io/nats.go"
)

const (
	subjectPrizeCreated    = "prize.created"
	subjectLaureateCreated = "laureate.created"
)

type Subscriber struct {
	conn *nats.Conn

	subs []*nats.Subscription
}

func New(conn *nats.Conn) *Subscriber {
	return &Subscriber{conn: conn}
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
