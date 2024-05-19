package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	UserDBPassword string `env:"UserDatabasePassword"`
	UserDBName     string `env:"UserDatabaseName"`
	DBName         string `env:"DatabaseName"`
	DriverDBName   string `env:"DriverDatabaseName"`
	URL            string `env:"URL"`
	ApiKey         string `env:"ApiKey"`
	Email          string `env:"Email"`
	EmailPassword  string `env:"EmailPassword"`
}

func LoadENV(filename string) *Config {
	err := godotenv.Load(filename)
	if err != nil {
		log.Panic().Err(err).Msg(" does not load .env")
	}
	log.Info().Msg("successfully load .env")
	cfg := Config{}
	return &cfg

}

func (cfg *Config) ParseENV() {

	err := env.Parse(cfg)
	if err != nil {
		log.Panic().Err(err).Msg(" unable to parse environment variables")
	}
	log.Info().Msg("successfully parsed .env")
}
