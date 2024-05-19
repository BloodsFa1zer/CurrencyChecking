package database

import (
	"errors"
)

type MockDb struct {
	EmailExists bool
	InsertError error
}

func (m *MockDb) InsertUser(user User) error {
	if m.EmailExists {
		return errors.New("such Email already exists")
	}
	if m.InsertError != nil {
		return m.InsertError
	}
	return nil
}

func (m *MockDb) SelectUsersEmail() ([]string, error) {
	return []string{}, nil
}
