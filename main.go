package main

import (
	"CurrencyChecking/communication"
	"CurrencyChecking/database"
	"CurrencyChecking/routes"
	"CurrencyChecking/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func main() {
	e := echo.New()
	// e.HidePort = true

	go communication.NewUserEmailSender(service.NewUserService(database.NewUserDatabase(), validate),
		database.NewUserDatabase()).ScheduleEmailSender()

	routes.UserRoute(e)
	database.NewUserDatabase()

	e.Logger.Fatal(e.Start(":6000"))
}
