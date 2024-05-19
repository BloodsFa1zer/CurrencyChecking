package database

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type User struct {
	ID        int64  `db:"ID"`
	Email     string `db:"email" json:"Email" validate:"required"`
	CreatedAt string `db:"created_at" json:"CreatedAt"`
}

type UserDatabase struct {
	Connection *sql.DB
}

func (db *UserDatabase) InsertUser(user User) error {
	formattedTime := time.Now().Format("2006.01.02 15:04")

	sqlInsert := "INSERT INTO Users (email, created_at) VALUES ($1, $2);"

	var emailExists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)", user.Email).Scan(&emailExists)
	if err != nil {
		return err
	}

	if emailExists {
		return errors.New("such Email already exists")
	}

	_, err = db.Connection.Exec(sqlInsert, user.Email, formattedTime)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserDatabase) SelectUsersEmail() ([]string, error) {
	rows, err := db.Connection.Query("SELECT email FROM users")
	if err != nil {
		log.Warn().Err(err).Msg("can`t select user`s email")
		return nil, err
	}

	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Warn().Err(err).Msg("can`t scan user`s email")
			return nil, err
		}
		emails = append(emails, email)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}
