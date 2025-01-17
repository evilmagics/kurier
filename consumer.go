package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	engine
}

func (cons *Consumer) consume(ds <-chan amqp.Delivery, fn ConsumeFunc) {
	for {
		log.Debug().Msg("start listening to consume")

		select {
		case d, ok := <-ds:
			if !ok {
				return
			}
			log.Info().Str("body", string(d.Body)).Msg("consume")

			// Do something on consumtions
			if fn != nil {
				fn(d)
			}

			d.Ack(false)
		}
	}
}

func (cons *Consumer) Listen(config ConsumerConfig) error {
	// Set our quality of service.  Since we're sharing 3 consumers on the same
	// channel, we want at least 2 messages in flight.
	err := cons.Chann.Qos(2, 0, false)
	if err != nil {
		return err
	}

	published, err := cons.Chann.Consume(
		config.Queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return err
	}

	go cons.consume(published, config.HandleConsume)

	return nil
}

func NewConsumer(config Config) (*Consumer, error) {
	cons := &Consumer{
		engine: createEngine(config),
	}

	if err := cons.load(config); err != nil {
		return nil, err
	}

	return cons, nil
}
