package database

type DbInterface interface {
	InsertUser(user User) error
	SelectUsersEmail() ([]string, error)
}
