package database

import (
	"database/sql"
	"fmt"

	"CurrencyChecking/config"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/zerolog/log"
)

func NewUserDatabase() *UserDatabase {
	cfg := config.LoadENV(".env")

	connStr := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBPort, cfg.UserDBName, cfg.UserDBPassword, cfg.DBName)

	log.Info().Msgf("Connection string: %s", connStr)

	db, err := sql.Open(cfg.DriverDBName, connStr)
	if err != nil {
		log.Warn().Err(err).Msg("Unable to connect to database")
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to ping the database")
		return nil
	}
	log.Info().Msg("Successfully connected to the database.")

	return &UserDatabase{Connection: db}
}
