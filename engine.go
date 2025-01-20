package kurier

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type engine struct {
	URL        string
	Conn       *amqp.Connection
	Chann      *amqp.Channel
	closeChann chan *amqp.Error
	quitChann  chan bool
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
		return err
	}

	e.closeChann = make(chan *amqp.Error)
	e.Conn.NotifyClose(e.closeChann)

	return nil
}

func (e *engine) Shutdown() {
	e.quitChann <- true
	log.Info().Msg("shutting down rabbitMQ connection")
	<-e.quitChann
}

func (e *engine) handleDisconnect(config Config) {
	for {
		select {
		case errChann := <-e.closeChann:
			if errChann != nil {
				log.Error().Err(errChann).Msg("rabbitMQ disconnected")
			}
		case <-e.quitChann:
			e.Conn.Close()
			log.Info().Msg("rabbitMQ has been shut down")
			e.quitChann <- true
			return
		}
		log.Info().Msg("trying to reconnect to rabbitMQ")

		time.Sleep(config.RetryConnInterval)

		if err := e.connect(); err != nil {
			log.Error().Err(err).Msg("reconnecting rabbitMQ failed")
		}
	}
}

func createEngine(config Config) engine {
	return engine{
		URL: config.AMQPUrl,
	}
}
