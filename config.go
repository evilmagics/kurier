package kurier

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type AMQP struct {
	URL      string `default:"amqp://guest:guest@127.0.0.1:5672/guest"`
	Exchange string `default:"amq.direct"`
}

type Config struct {
	AppName           string           `json:"app_name" default:"rabbitmq"`
	AppVersion        string           `json:"app_version" default:"1.0.0"`
	AMQPUrl           string           `json:"amqp_url"`
	Exchange          []ExchangeConfig `json:"exchanges"`
	Queue             []QueueConfig    `json:"queue"`
	LogLevel          zerolog.Level    `json:"-"`
	logger            zerolog.Logger   `json:"-"`
	RetryConnInterval time.Duration    `json:"-"`
}

type ExchangeConfig struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	Bind       QueueBindConfig
}

type QueueBindConfig struct {
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}

// queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table
type ConsumeFunc func(d amqp.Delivery) error
type ConsumerConfig struct {
	Queue         string
	Consumer      string
	AutoAck       bool
	Exclusive     bool
	NoLocal       bool
	NoWait        bool
	Args          amqp.Table
	HandleConsume ConsumeFunc
	Workers       int
	EnableMetrics bool
}

func DefaultExchange(name string) ExchangeConfig {
	return ExchangeConfig{
		Name:       name,
		Kind:       "x-delayed-message",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args: amqp.Table{
			"x-delayed-type": "direct",
		},
	}
}

func DefaultQueue(name string, bind QueueBindConfig) QueueConfig {
	return QueueConfig{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Bind:       bind,
	}
}

func DefaultQueueBind(key, exchange string) QueueBindConfig {
	return QueueBindConfig{
		Key:      key,
		Exchange: exchange,
		NoWait:   false,
	}
}

func DefaultConsumer(queue, consumer string) ConsumerConfig {
	return ConsumerConfig{
		Queue:     queue,
		Consumer:  consumer,
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
	}
}
