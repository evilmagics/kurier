package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectFunc func() error

func createConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	chann, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	log.Info().Msg("connection to rabbitMQ established")

	return conn, chann, nil
}
