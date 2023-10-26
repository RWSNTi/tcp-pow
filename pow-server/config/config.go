package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

// Config - App configuration
var Config struct {
	ServerHost                 string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort                 int           `env:"SERVER_PORT" envDefault:"4242"`
	Difficulty                 int           `env:"DIFFICULTY" envDefault:"4"`               // number of target values in hash prefix, required in hash to pass pow challenge
	TargetValue                string        `env:"TARGET_VALUE" envDefault:"0"`             // target value for result hash, usually "0"
	ClientConnectionTimeoutInt int           `env:"TIMEOUT" envDefault:"5"`                  // value of client timeout in seconds
	QuotesApiUrl               string                                                        // address of quotes API resource
	ClientConnectionTimeout    time.Duration
}

// LoadConfig from env and config file
func LoadConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}

	err = applyConfigFile("../config.toml")
	if err != nil {
		return err
	}

	Config.ClientConnectionTimeout = time.Second * time.Duration(Config.ClientConnectionTimeoutInt)

	log.Println("config loaded")
	return nil
}

func applyConfigFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("config file %s not found", filePath)
	}

	_, err := toml.DecodeFile(filePath, &Config)
	if err != nil {
		return fmt.Errorf("decode config file %s error: %v", filePath, err)
	}

	return nil
}
