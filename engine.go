package kurier

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type engine struct {
	URL        string
	Conn       *amqp.Connection
	Chann      *amqp.Channel
	closeChann chan *amqp.Error
	quitChann  chan bool
	logger     zerolog.Logger
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

func (e *engine) connect() (err error) {
	e.Conn, e.Chann, err = createConnection(e.URL)
	if err != nil {
		e.logger.Error().Err(err).Msg("failed to connect to rabbitMQ")
		return err
	}

	e.logger.Info().Msg("Connected to RabbitMQ")

	e.closeChann = make(chan *amqp.Error)
	e.Conn.NotifyClose(e.closeChann)

	return nil
}

func (e *engine) Shutdown() {
	e.quitChann <- true
	e.logger.Info().Msg("shutting down rabbitMQ connection")
	<-e.quitChann
}

func (e *engine) handleDisconnect(config Config) {
	for {
		select {
		case errChann := <-e.closeChann:
			if errChann != nil {
				e.logger.Error().Err(errChann).Msg("rabbitMQ disconnected")
			}
		case <-e.quitChann:
			e.Conn.Close()
			e.logger.Info().Msg("rabbitMQ has been shut down")
			e.quitChann <- true
			return
		}
		e.logger.Info().Msg("trying to reconnect to rabbitMQ")

		time.Sleep(config.RetryConnInterval)

		if err := e.connect(); err != nil {
			e.logger.Error().Err(err).Msg("reconnecting rabbitMQ failed")
		}
	}
}

func createEngine(config Config) engine {
	if config.Logger == nil {
		config.Logger = &log
	}
	return engine{
		URL:    config.AMQPUrl,
		logger: *config.Logger,
	}
}
