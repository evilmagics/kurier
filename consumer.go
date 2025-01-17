package main

import (
	"net/http"

	"github.com/alitto/pond/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	engine
	pool pond.Pool
}

func (cons *Consumer) registerPromotheus(pool pond.Pool) {
	// Worker pool metrics
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_workers_running",
			Help: "Number of running worker goroutines",
		},
		func() float64 {
			return float64(pool.RunningWorkers())
		}))

	// Task metrics
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_submitted_total",
			Help: "Number of tasks submitted",
		},
		func() float64 {
			return float64(pool.SubmittedTasks())
		}))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "pool_tasks_waiting_total",
			Help: "Number of tasks waiting in the queue",
		},
		func() float64 {
			return float64(pool.WaitingTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_successful_total",
			Help: "Number of tasks that completed successfully",
		},
		func() float64 {
			return float64(pool.SuccessfulTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_failed_total",
			Help: "Number of tasks that completed with panic",
		},
		func() float64 {
			return float64(pool.FailedTasks())
		}))
	prometheus.MustRegister(prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "pool_tasks_completed_total",
			Help: "Number of tasks that completed either successfully or with panic",
		},
		func() float64 {
			return float64(pool.CompletedTasks())
		}))

	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())
}

func (cons *Consumer) consume(config ConsumerConfig) {
	ds, err := cons.Chann.Consume(
		config.Queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed create delivery consumer")
		return
	}

	for {
		select {
		case d, ok := <-ds:
			if !ok {
				return
			}
			cons.consumeJob(d, config.HandleConsume)

			d.Ack(false)
		}
	}
}

func (cons *Consumer) consumeJob(d amqp.Delivery, fn ConsumeFunc) {
	cons.pool.Submit(func() {
		log.Info().Str("body", string(d.Body)).Msg("consume")

		// Do something on consumtions
		if fn != nil {
			fn(d)
		}
	})
}

func (cons *Consumer) createWorkers(count int) {
	// Create minimum workers
	if count <= 0 {
		count = 1
	}
	cons.pool = pond.NewPool(count)
}

func (cons *Consumer) Listen(config ConsumerConfig) error {
	// Set our quality of service.  Since we're sharing 3 consumers on the same
	// channel, we want at least 2 messages in flight.
	err := cons.Chann.Qos(config.Workers, 0, false)
	if err != nil {
		return err
	}

	// Create workers
	cons.createWorkers(config.Workers)

	go cons.consume(config)

	// Start exposing Prometheus metrics
	if config.EnableMetrics {
		cons.registerPromotheus(cons.pool)

		log.Info().Msg("Prometheus metrics server started on localhost:8989")
		if err := http.ListenAndServe(":8989", nil); err != nil {
			log.Warn().Err(err).Msg("Failed to start Prometheus metrics server")
		}
	}

	return nil
}

func NewConsumer(config Config) (*Consumer, error) {
	cons := &Consumer{
		engine: createEngine(config),
	}

	if err := cons.load(config); err != nil {
		return nil, err
	}

	return cons, nil
}
