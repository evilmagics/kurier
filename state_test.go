package kurier

import (
	"fmt"
	"testing"
	"time"
)

func TestNetworkState(t *testing.T) {
	config := Config{
		AppName:           "test_consumer",
		RetryConnInterval: 5 * time.Second,
		AMQPUrl:           "amqp://guest:guest@localhost:5672",
		// NetworkStateHook:  hook,
		NetworkStateHook: WatchState,
	}

	cons, err := NewConsumer(config)
	if err != nil {
		t.Fatal("Failed create new RabbitMQ producer!", err.Error())
	}
	// go WatchState(cons.NetworkState())

	defer cons.Shutdown()

	err = cons.Listen(DefaultConsumer("payments-status", "payments.status.consumer"))
	if err != nil {
		t.Fatal("Failed listening RabbitMQ consumer!", err.Error())
	}

}

func WatchState(state NetworkState) {
	switch state {
	case NetworkConnecting:
		fmt.Println("Connecting")
	case NetworkConnectingFailed:
		fmt.Println("Connecting Failed")
	case NetworkConnected:
		fmt.Println("Connected")
	case NetworkDisconnecting:
		fmt.Println("Disconnecting")
	case NetworkDisconnectingFailed:
		fmt.Println("Disconnecting Failed")
	case NetworkDisconnected:
		fmt.Println("Disconnected")
	case NetworkReconnecting:
		fmt.Println("Reconnecting")
	case NetworkReconnectingFailed:
		fmt.Println("Reconnecting Failed")
	case NetworkReconnected:
		fmt.Println("Reconnected")
	case NetworkShutdown:
		fmt.Println("Shutdown")
	default:
		fmt.Println("Unknown State")
	}

}
