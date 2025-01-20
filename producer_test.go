package kurier

import (
	"encoding/json"
	"testing"
	"time"
)

func testProducerConfig() Config {
	return Config{
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
}

func TestProducer(t *testing.T) {
	config := testProducerConfig()

	prod, err := NewProducer(config)
	if err != nil {
		t.Fatal("Failed create new RabbitMQ producer!")
	}

	defer prod.Shutdown()

	data := map[string]interface{}{
		"user_id":    "usr_0001",
		"billing_id": "trx_0001",
		"created_at": time.Now().Format(time.DateTime),
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed marshal body payload! %s", err.Error())
	}

	err = prod.PublishDelayed(config.Exchange[0].Name, config.Queue[0].Bind.Key, body, 3*time.Second)
	if err != nil {
		t.Fatalf("failed to publish into rabbitmq: %v", err)
	}

	t.Fail()
}
func BenchmarkProducer(t *testing.B) {
	config := testProducerConfig()

	prod, err := NewProducer(config)
	if err != nil {
		t.Fatal("Failed create new RabbitMQ producer!")
	}

	defer prod.Shutdown()

	data := map[string]interface{}{
		"user_id":    "usr_0001",
		"billing_id": "trx_0001",
		"created_at": time.Now().Format(time.DateTime),
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed marshal body payload! %s", err.Error())
	}

	err = prod.PublishDelayed(config.Exchange[0].Name, config.Queue[0].Bind.Key, body, 2*time.Second)
	if err != nil {
		t.Fatalf("failed to publish into rabbitmq: %v", err)
	}
}
