package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/mocks"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	auth := mocks.NewMockAuthService(ctrl)
	cfg := config.GetConfig()
	handler := NewHandler(nil, auth, &cfg)
	httpSrv := httptest.NewServer(handler.Router)

	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		loginDTO      models.LoginUserDTO
		tokenDTO      *models.ReadTokenDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success login",
			expectedCode:  http.StatusOK,
			loginDTO:      models.LoginUserDTO{UserName: "test_user", Password: "test123!"},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Error wrong password",
			expectedCode:  http.StatusUnauthorized,
			loginDTO:      models.LoginUserDTO{UserName: "test_user", Password: "test123!"},
			tokenDTO:      nil,
			serviceError:  service.ErrWrongPassword,
			serviceCalled: true,
		},
		{
			name:          "Error user not found",
			expectedCode:  http.StatusUnauthorized,
			loginDTO:      models.LoginUserDTO{UserName: "test_user", Password: "test123!"},
			tokenDTO:      nil,
			serviceError:  service.ErrUserNotFound,
			serviceCalled: true,
		},
		{
			name:          "User name validation error",
			expectedCode:  http.StatusUnprocessableEntity,
			loginDTO:      models.LoginUserDTO{UserName: "user", Password: "test123!"},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "Password validation error",
			expectedCode:  http.StatusUnprocessableEntity,
			loginDTO:      models.LoginUserDTO{UserName: "test_user", Password: "test"},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				auth.EXPECT().Login(tc.loginDTO).Return(tc.tokenDTO, tc.serviceError)
			}
			body, _ := json.Marshal(tc.loginDTO)
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = httpSrv.URL + "/v1.0/auth/login"
			resp, err := req.SetBody(body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	auth := mocks.NewMockAuthService(ctrl)
	cfg := config.GetConfig()
	handler := NewHandler(nil, auth, &cfg)
	httpSrv := httptest.NewServer(handler.Router)

	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		registerDTO   models.RegisterUserDTO
		tokenDTO      *models.ReadTokenDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:         "Success login",
			expectedCode: http.StatusCreated,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test_user",
				Password:        "test123!",
				PasswordConfirm: "test123!",
				FirstName:       "John",
				LastName:        "Doe",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:         "Password missmatch error",
			expectedCode: http.StatusUnprocessableEntity,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test_user",
				Password:        "test123!!!",
				PasswordConfirm: "test123!",
				FirstName:       "John",
				LastName:        "Doe",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:         "User name validation error",
			expectedCode: http.StatusUnprocessableEntity,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test",
				Password:        "test123!",
				PasswordConfirm: "test123!",
				FirstName:       "John",
				LastName:        "Doe",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:         "Password validation error",
			expectedCode: http.StatusUnprocessableEntity,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test_user",
				Password:        "test",
				PasswordConfirm: "test",
				FirstName:       "John",
				LastName:        "Doe",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:         "First name validation error",
			expectedCode: http.StatusUnprocessableEntity,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test_user",
				Password:        "test123!",
				PasswordConfirm: "test123!",
				FirstName:       "",
				LastName:        "Doe",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:         "Last name validation error",
			expectedCode: http.StatusUnprocessableEntity,
			registerDTO: models.RegisterUserDTO{
				UserName:        "test_user",
				Password:        "test123!",
				PasswordConfirm: "test123!",
				FirstName:       "John",
				LastName:        "",
			},
			tokenDTO:      &models.ReadTokenDTO{AccessToken: "test", RefreshToken: "test"},
			serviceError:  nil,
			serviceCalled: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				auth.EXPECT().Register(tc.registerDTO).Return(tc.tokenDTO, tc.serviceError)
			}
			body, _ := json.Marshal(tc.registerDTO)
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = httpSrv.URL + "/v1.0/auth/register"
			resp, err := req.SetBody(body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
