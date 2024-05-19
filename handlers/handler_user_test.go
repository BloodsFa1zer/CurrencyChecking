package handlers

import (
	"CurrencyChecking/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestContext creates a new Echo context and a response recorder for testing purposes.
func newTestContext(method, path, request string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()

	// Create a request
	req := httptest.NewRequest(method, path, strings.NewReader(request))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	return ctx, rec
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockUserServiceInterface(ctrl)

	handler := NewUserHandler(mockService)
	tt := []struct {
		scenario           string
		request            string
		expectedResponse   string
		expectedStatusCode int
		expectedError      error
	}{
		{scenario: "success",
			request: `{"Email":"test@email.com"}`,
			expectedResponse: `{
									"status": 200,
									"message": "success",
									"data": {
										"data": "Email added"
									}
								}`,
			expectedStatusCode: http.StatusOK,
			expectedError:      nil,
		},
		{scenario: "email_validation_failed_domain",
			request: `{"Email":"testEmail.com"}`,
			expectedResponse: `{
									"status": 400,
									"message": "error",
									"data": {
										"data": "local or domain part cannot be empty in the email"
									}
								}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      errors.New("local or domain part cannot be empty in the email"),
		},
		{scenario: "email_validation_failed_void",
			request: `{"Email":""}`,
			expectedResponse: `{
									"status": 400,
									"message": "error",
									"data": {
										"data": "email must contain a single '@' character"
									}
								}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      errors.New("email must contain a single '@' character"),
		},
		{scenario: "email_validation_failed_irrelevant_characters",
			request: `{"Email":"test..test@example.com"}`,
			expectedResponse: `{
									"status": 400,
									"message": "error",
									"data": {
										"data": "consecutive special characters '.' are not allowed in the local part"
									}
								}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      errors.New("consecutive special characters '.' are not allowed in the local part"),
		},
		{scenario: "email_exists",
			request: `{"Email":"test@email.com"}`,
			expectedResponse: `{
									"status": 409,
									"message": "error",
									"data": {
										"data": "such Email already exists"
									}
								}`,
			expectedStatusCode: http.StatusConflict,
			expectedError:      errors.New("such Email already exists"),
		},
		{scenario: "database_error",
			request: `{"Email":"test@email.com"}`,
			expectedResponse: `{
									"status": 400,
									"message": "error",
									"data": {
										"data": "can't insert user"
									}
								}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      errors.New("can't insert user"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.scenario, func(t *testing.T) {
			mockService.EXPECT().CreateUser(gomock.Any()).Return(tc.expectedError, tc.expectedStatusCode).Times(1)

			ctx, rec := newTestContext(http.MethodPost, "/subscribe", tc.request)

			err := handler.CreateUser(ctx)

			// Check the response
			if assert.NoError(t, err) {
				t.Log("Actual Response:", rec.Body.String())

				assert.Equal(t, tc.expectedStatusCode, rec.Code)

				assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			}
		})
	}
}
