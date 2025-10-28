package kurier

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type NetworkState int

const (
	NetworkConnecting NetworkState = iota
	NetworkConnectingFailed
	NetworkConnected
	NetworkDisconnecting
	NetworkDisconnectingFailed
	NetworkDisconnected
	NetworkReconnecting
	NetworkReconnectingFailed
	NetworkReconnected
	NetworkShutdown
)

type NetworkStateHook func(state NetworkState)

type engine struct {
	URL         string
	Conn        *amqp.Connection
	Chann       *amqp.Channel
	closeChann  chan *amqp.Error
	quitChann   chan bool
	logger      zerolog.Logger
	networkHook NetworkStateHook
}

func (e *engine) load(config Config) error {
	var err error

	if err = e.connect(); err != nil {
		return err
	}

	if err = declareExchange(e.Chann, config.Exchange); err != nil {
		return err
	}

	// declare queue
	if err = declareQueue(e.Chann, config.Queue); err != nil {
		return err
	}

	e.quitChann = make(chan bool)

	go e.handleDisconnect(config)

	return nil
}
func (e *engine) setNetworkState(state NetworkState) {

	if e.networkHook != nil {
		e.networkHook(state)
	}
}

func (e *engine) connect() (err error) {
	e.Conn, e.Chann, err = createConnection(e.URL)
	if err != nil {
		e.setNetworkState(NetworkDisconnectingFailed)
		return err
	}

	e.setNetworkState(NetworkConnected)

	e.closeChann = make(chan *amqp.Error)
	e.Conn.NotifyClose(e.closeChann)

	return nil
}

func (e *engine) Shutdown() {
	e.quitChann <- true
	e.setNetworkState(NetworkShutdown)
	<-e.quitChann
}

func (e *engine) handleDisconnect(config Config) {
	for {
		select {
		case errChann := <-e.closeChann:
			if errChann != nil {
				e.setNetworkState(NetworkDisconnected)
			}
		case <-e.quitChann:
			e.Conn.Close()
			e.setNetworkState(NetworkShutdown)
			e.quitChann <- true
			return
		}

		e.setNetworkState(NetworkReconnecting)

		time.Sleep(config.RetryConnInterval)

		if err := e.connect(); err != nil {
			e.setNetworkState(NetworkDisconnectingFailed)
		}
	}
}

func createEngine(config Config) engine {
	if config.Logger == nil {
		config.Logger = &log
	}
	e := engine{
		URL:         config.AMQPUrl,
		logger:      *config.Logger,
		networkHook: config.NetworkStateHook,
	}

	return e
}
