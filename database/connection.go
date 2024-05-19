package database

import (
	"CurrencyChecking/config"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func NewUserDatabase() *UserDatabase {
	cfg := config.LoadENV("config/.env")
	cfg.ParseENV()

	connStr := "user=" + cfg.UserDBName + " password=" + cfg.UserDBPassword + " dbname=" + cfg.DBName + " sslmode=disable"
	db, err := sql.Open(cfg.DriverDBName, connStr)
	if err != nil {
		log.Warn().Err(err).Msg("can`t connect to database")
	}

	err = db.Ping()
	if err != nil {
		log.Warn().Err(err).Msg("failed to ping the database")
		return nil
	}
	log.Info().Msg("successfully connected to the database.")

	return &UserDatabase{Connection: db}
}
