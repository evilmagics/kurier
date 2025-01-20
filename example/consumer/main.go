package main

import (
	"os"
	"time"

	kurier "github.com/evilmagics/kurier"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

var log zerolog.Logger = zerolog.
	New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
	With().
	Timestamp().
	Logger()

func handleConsumption(d amqp.Delivery) error {
	log.Debug().
		Str("process_timestamp", d.Timestamp.Format(time.DateTime)).
		Str("body", string(d.Body)).
		Msg("Consume")

	return nil
}
func main() {
	config := kurier.Config{
		AppName:           "test_consumer",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
	}
	consumerConfig := kurier.ConsumerConfig{
		Queue:         "payments-status",
		Consumer:      "payments.status.consumer",
		AutoAck:       false,
		Exclusive:     false,
		NoLocal:       false,
		NoWait:        false,
		HandleConsume: handleConsumption,
		Workers:       10,
		EnableMetrics: true,
	}

	cons, err := kurier.NewConsumer(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed create new RabbitMQ producer!")
	}

	defer cons.Shutdown()

	err = cons.Listen(consumerConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed listening RabbitMQ consumer!")
	}
}
