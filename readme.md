
<img style="display: block; margin-left: auto; margin-right: auto; width: 30%;" src="assets/logo.png" alt="Kurier Logo">

<div style="text-align:center">
Handle <b><i>RabbitMQ</i></b> delayed message and task for Go.
</div>

# ⚡ Kurier

The project seems to be designed to handle messages that need to be processed or delivered after a specified delay. This could be useful in various scenarios such as scheduled notifications, timed events, or implementing retry mechanisms.

## 🔔 Features

- **Message Scheduling**: Ability to schedule messages for future processing or delivery.
- **Configurable Delays**: Users can specify custom delay durations for each message.
- **Persistent Storage**: Messages are stored safely to survive application restarts.
- **Scalable Architecture**: Designed to handle a large volume of delayed messages efficiently.
- **Producer-Consumer Model**: Separate components for message production and consumption.

## 🧩 Components

1. **Producer**: Responsible for creating and scheduling delayed messages.
2. **Consumer**: Processes or delivers messages when their delay time has elapsed.
3. **Engine**: Core component managing the delay mechanism and message flow.
4. **Config**: Allows customization of application behavior through configuration files.
5. **Logger**: Provides detailed logging for monitoring and debugging.
6. **Connections**: Manages any necessary RabbitMQ connections.

## 🪄 Use Cases

- Sending reminder emails or notifications at specific times.
- Implementing retry logic for failed operations.
- Scheduling tasks or jobs to run at future times.
- Managing time-sensitive workflows in distributed systems.

## ⚙️ Installation

  1. **Prerequisites**
     - RabbitMQ 4.0.5+
     - Erlang 27.2+
     - Go 1.16+

  2. **Installation**

```bash
go get github.com/evilmagics/kurier
```

**Recommendation**: Clean and update package before using.
	
```bash
go mod tidy
```

## 🎯 Usage

Here's a simple example of how to use the Kurier library in your application. 

### Producer

The producer has two functions that both send data to the exchange on the RabbitMQ service. For configuration can read on [configuration](#configuration) section.

Make sure the RabbitMQ service is running on your platform to receive data from the producer. Receiving data that has been sent by the producer can be seen in the [consumer](#consumer) section.

1. **Initialization**

```go
package main 

func main() {
    config := Config{
        AppName:           "test_producer",
        AppVersion:        "1.0.0",
        RetryConnInterval: 5 * time.Second,
        AMQPUrl:           "amqp://guest:guest@localhost:5672",
        Exchange: []ExchangeConfig{
            DefaultExchange("payments"),
        },
        Queue: []QueueConfig{
            DefaultQueue(
                "payments-status",
                DefaultQueueBind("payment.event.status", "payments"),
            ),
        },
    }

    prod, err := NewProducer(config)
    if err != nil {
        t.Fatal("Failed create new RabbitMQ producer!")
    }

    defer prod.Shutdown()
}
```

2. **Preparing body payload**

```go
data := map[string]interface{}{
	"user_id":    "usr_0001",
}

body, err := json.Marshal(data)
if err != nil {
	t.Fatalf("Failed marshal body payload! %s", err.Error())
}
```

3. **Publish**

    Publish to exchanges directly without any delay.

```go
err = prod.Publish(config.Exchange[0].Name, config.Queue[0].Bind.Key, body)
if err != nil {
	t.Fatalf("failed to publish into rabbitmq: %v", err)
}
```

4. **Publish Delayed**

Publish to exchanges using delay on milliseconds (ms) refer to [limitation](https://github.com/rabbitmq/rabbitmq-delayed-message-exchange?tab=readme-ov-file#limitations) and should using `time.Duration` from build-in golang time package.

```go
err = prod.PublishDelayed(config.Exchange[0].Name, config.Queue[0].Bind.Key, body, 5*time.Second)
if err != nil {
	t.Fatalf("failed to publish into rabbitmq: %v", err)
}
```

### Consumer

The consumer is responsible for receiving and processing messages from RabbitMQ queues. Here's how to set up and use a consumer:

1. **Initialization**

   Consumer initialization does not require as much configuration as producer. The configuration for consumer requires RabbitMQ AMQP URL to connect to RabbitMQ.

```go
package main

import (
	"fmt"
	"time"
	"github.com/evilmagics/kurier"
)

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
}
```

2. **Consumption**

```gol
err = cons.Listen(DefaultConsumer("payments-status", "payments.status.consumer"))
if err != nil {
	log.Fatal().Err(err).Msg("Failed listening RabbitMQ consumer!")
}
```

## Configuration

Setting up the Config struct with necessary RabbitMQ connection details and queue configurations.

### General Configuration

```go
type Config struct {
	AppName           string           `json:"app_name" default:"rabbitmq"`
	AppVersion        string           `json:"app_version" default:"1.0.0"`
	AMQPUrl           string           `json:"amqp_url"`
	Exchange          []ExchangeConfig `json:"exchanges"`
	Queue             []QueueConfig    `json:"queue"`
	LogLevel          zerolog.Level    `json:"-"`
	RetryConnInterval time.Duration    `json:"-"`
}
```

### Exchange Configuration

```go
type ExchangeConfig struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}
```

### Queue Configuration

```go
type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	Bind       QueueBindConfig
}
```

### Queue Binding Configuration

```go
type QueueBindConfig struct {
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}
```

### Consume Function

This function type is used to define how each message should be processed when it's consumed from a queue.

```go
type ConsumeFunc func(d amqp.Delivery) error
```

### Consumer Listener Configuration

```go
type ConsumerConfig struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
	onConsume ConsumeFunc
}
```

## 💡 Acknowledgements

Inspired by the [RabbitMQ Delayed Message Exchange Plugin](https://github.com/rabbitmq/rabbitmq-delayed-message-exchange) and [Ghith's RabbitMQ Delayed Example](https://github.com/ghigt/rabbitmq-delayed).

We're committed to improving and expanding this library. Your feedback and contributions are welcome!
