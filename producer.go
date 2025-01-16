package main

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	engine
}

func (prod *Producer) Publish(exchange, key string, body []byte) error {
	return prod.publish(exchange, key, body, 0)
}
func (prod *Producer) PublishDelayed(exchange, key string, body []byte, delay time.Duration) error {
	return prod.publish(exchange, key, body, delay)
}
func (prod *Producer) publish(exchange, key string, body []byte, delay time.Duration) error {
	headers := make(amqp.Table)

	log.Debug().Str("exchange", exchange).Str("key", key).Msg("publishing")

	if delay > 0 {
		headers["x-delay"] = delay.Milliseconds()
	}

	return prod.Chann.Publish(exchange, key, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         body,
		Headers:      headers,
	})
}

func NewProducer(config Config) (*Producer, error) {
	prod := &Producer{
		engine: createEngine(config),
	}

	if err := prod.load(config); err != nil {
		return nil, err
	}

	return prod, nil
}
