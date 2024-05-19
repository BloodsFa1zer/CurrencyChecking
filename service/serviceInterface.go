package service

import "CurrencyChecking/database"

type UserServiceInterface interface {
	CreateUser(user database.User) (error, int)
	GetRate() (string, error, int)
}
