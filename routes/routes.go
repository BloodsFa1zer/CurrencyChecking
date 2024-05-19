package routes

import (
	"CurrencyChecking/database"
	"CurrencyChecking/handlers"
	"CurrencyChecking/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()
var userHandler = handlers.NewUserHandler(service.NewUserService(database.NewUserDatabase(), validate))

func UserRoute(e *echo.Echo) {
	e.GET("/rate", userHandler.GetRate)
	e.POST("/subscribe", userHandler.CreateUser)
}
