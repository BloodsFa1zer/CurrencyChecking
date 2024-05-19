package handlers

import (
	"CurrencyChecking/database"
	"CurrencyChecking/response"
	"CurrencyChecking/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: service}
}

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var user database.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus := uh.userService.CreateUser(user)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": "Email added"}})
}

func (uh *UserHandler) GetRate(c echo.Context) error {
	currencyRate, err, respStatus := uh.userService.GetRate()
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "Invalid status value"}})
	}

	return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "success", Data: &echo.Map{"data": currencyRate}})
}
