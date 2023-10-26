package config

import (
	"github.com/caarlos0/env/v6"
	"time"
)

// Config - App configuration
var Config struct {
	ServerHost              string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort              int           `env:"SERVER_PORT" envDefault:"4242"`
	IntervalInt             int           `env:"INTERVAL_SECONDS" envDefault:"5"`     // interval between reconnection to get new quote
	MaxSolveIterationLimit  int           `env:"ITERATIONS_LIMIT" envDefault:"50000"` // max iterations to prevent stuck on hard hashes (only for client)
	RequestIntervalDuration time.Duration                                             // time between restarts of problem solving func
}

// LoadConfig - Load configuration from env
func LoadConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	Config.RequestIntervalDuration = time.Second * time.Duration(Config.IntervalInt)

	return nil
}
