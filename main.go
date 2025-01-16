package main

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handleConsumption(d amqp.Delivery) {
	log.Debug().
		Str("process_timestamp", d.Timestamp.Format(time.DateTime)).
		Str("body", string(d.Body)).
		Msg("Consume")
}
func main() {
	config := Config{
		AppName:           "test_consumer",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
	}

	cons, err := NewConsumer(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed create new RabbitMQ producer!")
	}

	defer cons.Shutdown()

	err = cons.Listen(DefaultConsumer("payments-status", "payments.status.consumer"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed listening RabbitMQ consumer!")
	}

	for {
	}
}
