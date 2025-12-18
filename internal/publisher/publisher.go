package publisher

import (
	"encoding/json"
	"fmt"
	"ris/internal/domain"

	"github.com/nats-io/nats.go"
)

const (
	subjectPrizeCreated    = "prize.created"
	subjectLaureateCreated = "laureate.created"
)

type Publisher struct {
	broker *nats.Conn
}

func New(broker *nats.Conn) *Publisher {
	return &Publisher{broker: broker}
}

func (p *Publisher) PublishPrizeCreated(prize domain.Prize) error {
	data, err := json.Marshal(prize)
	if err != nil {
		return fmt.Errorf("error marshalling prize %v", err)
	}
	return p.broker.Publish(subjectPrizeCreated, data)
}

func (p *Publisher) PublishLaureateCreated(laureate domain.Laureate) error {
	data, err := json.Marshal(laureate)
	if err != nil {
		return fmt.Errorf("error marshalling laureate %v", err)
	}
	return p.broker.Publish(subjectLaureateCreated, data)
}
