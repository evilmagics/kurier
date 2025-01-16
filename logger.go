package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger = zerolog.
	New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
	With().
	Timestamp().
	Logger()

func initLogger(config *Config) {
	log = zerolog.
		New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
		Level(config.LogLevel).
		With().
		Timestamp().
		Logger()
}
