package kurier

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(chann *amqp.Channel, exchanges []ExchangeConfig) (err error) {
	for _, ex := range exchanges {
		err = chann.ExchangeDeclare(
			ex.Name,
			ex.Kind,
			ex.Durable,
			ex.AutoDelete,
			ex.Internal,
			ex.NoWait,
			ex.Args,
		)
		if err != nil {
			return errors.Wrapf(err, "declaring exchange %q", "delayed")
		}
	}
	return nil
}

func declareQueue(chann *amqp.Channel, queue []QueueConfig) error {
	for _, qc := range queue {
		q, err := chann.QueueDeclare(qc.Name, qc.Durable, qc.AutoDelete, qc.Exclusive, qc.NoWait, qc.Args)
		if err != nil {
			return err
		}
		err = chann.QueueBind(q.Name, qc.Bind.Key, qc.Bind.Exchange, qc.Bind.NoWait, qc.Bind.Args)
		if err != nil {
			return err
		}
	}

	return nil
}
