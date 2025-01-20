package kurier

import (
	"testing"
	"time"
)

func TestConsumption(t *testing.T) {
	config := Config{
		AppName:           "test_consumer",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
	}

	cons, err := NewConsumer(config)
	if err != nil {
		t.Fatal("Failed create new RabbitMQ producer!", err.Error())
	}

	defer cons.Shutdown()

	err = cons.Listen(DefaultConsumer("payments-status", "payments.status.consumer"))
	if err != nil {
		t.Fatal("Failed listening RabbitMQ consumer!", err.Error())
	}

	for {

	}
}
func BenchmarkConsumption(t *testing.B) {
	config := Config{
		AppName:           "test_consumer",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
	}

	cons, err := NewConsumer(config)
	if err != nil {
		t.Fatal("Failed create new RabbitMQ producer!", err.Error())
	}

	defer cons.Shutdown()

	err = cons.Listen(DefaultConsumer("payments-status", "payments.status.consumer"))
	if err != nil {
		t.Fatal("Failed listening RabbitMQ consumer!", err.Error())
	}

}
