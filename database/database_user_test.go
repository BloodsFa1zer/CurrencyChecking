package database

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	t.Run("Email does not exist, insert successful", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		userDB := &UserDatabase{Connection: db}

		user := User{
			Email: "test@gmail.com",
		}

		// Expect the SELECT query to check if the email exists
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)")).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		formattedTime := time.Now().Format("2006.01.02 15:04")

		// Expect the INSERT query to insert the new user
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Users (email, created_at) VALUES ($1, $2)")).
			WithArgs(user.Email, formattedTime).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Call the function being tested
		err = userDB.InsertUser(user)

		// Check for errors
		assert.NoError(t, err, "InsertUser should not return an error")

		// Verify that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Email already exists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		userDB := &UserDatabase{Connection: db}

		user := User{
			Email: "test@gmail.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)")).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err = userDB.InsertUser(user)

		assert.EqualError(t, err, "such Email already exists")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Database error during insert", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		userDB := &UserDatabase{Connection: db}

		user := User{
			Email: "test@gmail.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)")).
			WithArgs(user.Email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		formattedTime := time.Now().Format("2006.01.02 15:04")
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Users (email, created_at) VALUES ($1, $2)")).
			WithArgs(user.Email, formattedTime).
			WillReturnError(errors.New("insert failed"))

		err = userDB.InsertUser(user)

		assert.EqualError(t, err, "insert failed")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestSelectUsersEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userDB := &UserDatabase{
		Connection: db,
	}

	rows := sqlmock.NewRows([]string{"email"}).
		AddRow("test1@email.com").
		AddRow("test2@email.com")

	mock.ExpectQuery("SELECT email FROM users").WillReturnRows(rows)

	emails, err := userDB.SelectUsersEmail()
	assert.NoError(t, err, "SelectUsersEmail should not return an error")
	assert.Equal(t, []string{"test1@email.com", "test2@email.com"}, emails, "emails should match")

	// Verify the expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}

func TestInsertUser_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userDB := &UserDatabase{
		Connection: db,
	}

	// Mocking the query to check if email exists (return true)
	mock.ExpectQuery("SELECT EXISTS(.+)").WithArgs("test@email.com").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Mocking the actual insert query (should not be called in this test)
	user := User{
		Email: "test@email.com",
	}

	err = userDB.InsertUser(user)
	assert.Error(t, err, "InsertUser should return an error for existing email")
	assert.Equal(t, "such Email already exists", err.Error(), "error message should match")

	// Verify the expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "there were unfulfilled expectations")
}
