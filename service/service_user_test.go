package service

import (
	"CurrencyChecking/database"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"testing"

	"CurrencyChecking/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	var validate = validator.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDb := mocks.NewMockDbInterface(ctrl)
	userService := NewUserService(mockDb, validate)

	tt := []struct {
		scenario           string
		user               database.User
		mockBehavior       func()
		expectedError      error
		expectedStatusCode int
	}{
		{
			scenario: "success",
			user: database.User{
				Email: "test@email.com",
			},
			mockBehavior: func() {
				mockDb.EXPECT().InsertUser(gomock.Any()).Return(nil)
			},
			expectedError:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			scenario: "invalid_email",
			user: database.User{
				Email: "invalid-email",
			},
			mockBehavior: func() {
				// No behavior expected as it's an invalid email scenario
			},
			expectedError:      errors.New("email must contain a single '@' character"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario: "empty_email",
			user: database.User{
				Email: "",
			},
			mockBehavior: func() {
				// No behavior expected as it's an invalid email scenario
			},
			expectedError:      errors.New("email must contain a single '@' character"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario: "invalid_email",
			user: database.User{
				Email: "test..test@example.com",
			},
			mockBehavior: func() {
				// No behavior expected as it's an invalid email scenario
			},
			expectedError:      errors.New("consecutive special characters '.' are not allowed in the local part"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario: "email_exists",
			user: database.User{
				Email: "test@email.com",
			},
			mockBehavior: func() {
				mockDb.EXPECT().InsertUser(gomock.Any()).Return(errors.New("such Email already exists"))
			},
			expectedError:      errors.New("such Email already exists"),
			expectedStatusCode: http.StatusConflict,
		},
		{
			scenario: "db_error",
			user: database.User{
				Email: "test@email.com",
				// Add other required fields
			},
			mockBehavior: func() {
				mockDb.EXPECT().InsertUser(gomock.Any()).Return(errors.New("database error"))
			},
			expectedError:      errors.New("database error"),
			expectedStatusCode: http.StatusBadRequest, // Assuming bad request for database error
		},
	}

	for _, tc := range tt {
		t.Run(tc.scenario, func(t *testing.T) {
			// Define mock behavior for this test case
			tc.mockBehavior()

			err, statusCode := userService.CreateUser(tc.user)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatusCode, statusCode)
		})
	}
}
