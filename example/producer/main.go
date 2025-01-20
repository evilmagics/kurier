package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/evilmagics/kurier"
	"github.com/rs/zerolog"
)

var log zerolog.Logger = zerolog.
	New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
	With().
	Timestamp().
	Logger()

func main() {
	config := kurier.Config{
		AppName:           "test_producer",
		AppVersion:        "1.0.0",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
		Exchange: []kurier.ExchangeConfig{
			kurier.DefaultExchange("payments"),
		},
		Queue: []kurier.QueueConfig{
			kurier.DefaultQueue(
				"payments-status",
				kurier.DefaultQueueBind("payment.event.status", "payments"),
			),
		},
	}

	prod, err := kurier.NewProducer(config)
	if err != nil {
		log.Fatal().Msg("Failed create new RabbitMQ producer!")
	}

	defer prod.Shutdown()

	data := map[string]interface{}{
		"user_id":    "usr_0001",
		"billing_id": "trx_0001",
		"created_at": time.Now().Format(time.DateTime),
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed marshal body payload!")
	}

	err = prod.PublishDelayed(config.Exchange[0].Name, config.Queue[0].Bind.Key, body, 3*time.Second)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to publish into rabbitmq")
	}
}
