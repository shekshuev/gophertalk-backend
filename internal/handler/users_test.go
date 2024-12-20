package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/shekshuev/gophertalk-backend/internal/config"
	"github.com/shekshuev/gophertalk-backend/internal/mocks"
	"github.com/shekshuev/gophertalk-backend/internal/models"
	"github.com/shekshuev/gophertalk-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	users := mocks.NewMockUserService(ctrl)
	accessTokenSecret := "test"
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	accessToken, err := utils.CreateToken(
		cfg.AccessTokenSecret,
		"1",
		cfg.AccessTokenExpires,
	)
	assert.NoError(t, err, "error creating token")
	handler := NewHandler(users, nil, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		limit         string
		offset        string
		responseDTOs  []models.ReadUserDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success fetching users",
			expectedCode:  http.StatusOK,
			limit:         "10",
			offset:        "0",
			responseDTOs:  []models.ReadUserDTO{{ID: 1, UserName: "test_user"}},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid limit",
			expectedCode:  http.StatusOK,
			limit:         "-1",
			offset:        "0",
			responseDTOs:  []models.ReadUserDTO{{ID: 1, UserName: "test_user"}},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Service error",
			expectedCode:  http.StatusBadRequest,
			limit:         "10",
			offset:        "0",
			responseDTOs:  nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				users.EXPECT().GetAllUsers(gomock.Any(), gomock.Any()).Return(tc.responseDTOs, tc.serviceError)
			}

			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + "/v1.0/users?limit=" + tc.limit + "&offset=" + tc.offset
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	users := mocks.NewMockUserService(ctrl)
	accessTokenSecret := "test"
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	accessToken, err := utils.CreateToken(
		cfg.AccessTokenSecret,
		"1",
		cfg.AccessTokenExpires,
	)
	assert.NoError(t, err, "error creating token")
	handler := NewHandler(users, nil, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		userID        string
		responseDTO   *models.ReadUserDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success fetching user",
			expectedCode:  http.StatusOK,
			userID:        "1",
			responseDTO:   &models.ReadUserDTO{ID: 1, UserName: "test_user"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid user ID",
			expectedCode:  http.StatusNotFound,
			userID:        "abc",
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "User not found",
			expectedCode:  http.StatusNotFound,
			userID:        "2",
			responseDTO:   nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				users.EXPECT().GetUserByID(gomock.Any()).Return(tc.responseDTO, tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + "/v1.0/users/" + tc.userID
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_DeleteUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	users := mocks.NewMockUserService(ctrl)
	accessTokenSecret := "test"
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	accessToken, err := utils.CreateToken(
		cfg.AccessTokenSecret,
		"1",
		cfg.AccessTokenExpires,
	)
	assert.NoError(t, err, "error creating token")
	handler := NewHandler(users, nil, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		userID        string
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success deleting user",
			expectedCode:  http.StatusNoContent,
			userID:        "1",
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid user ID",
			expectedCode:  http.StatusNotFound,
			userID:        "abc",
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "User not found",
			expectedCode:  http.StatusNotFound,
			userID:        "2",
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				users.EXPECT().DeleteUser(gomock.Any()).Return(tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodDelete
			req.URL = httpSrv.URL + "/v1.0/users/" + tc.userID
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	users := mocks.NewMockUserService(ctrl)
	accessTokenSecret := "test"
	os.Setenv("ACCESS_TOKEN_SECRET", accessTokenSecret)
	os.Setenv("ACCESS_TOKEN_EXPIRES", "1h")
	cfg := config.GetConfig()
	accessToken, err := utils.CreateToken(
		cfg.AccessTokenSecret,
		"1",
		cfg.AccessTokenExpires,
	)
	assert.NoError(t, err, "error creating token")
	handler := NewHandler(users, nil, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		userID        string
		updateDTO     models.UpdateUserDTO
		responseDTO   *models.ReadUserDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success updating user",
			expectedCode:  http.StatusOK,
			userID:        "1",
			updateDTO:     models.UpdateUserDTO{FirstName: "Updated", LastName: "User"},
			responseDTO:   &models.ReadUserDTO{ID: 1, UserName: "updated_user"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid user ID",
			expectedCode:  http.StatusNotFound,
			userID:        "abc",
			updateDTO:     models.UpdateUserDTO{FirstName: "Updated", LastName: "User"},
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "Validation error",
			expectedCode:  http.StatusUnprocessableEntity,
			userID:        "1",
			updateDTO:     models.UpdateUserDTO{UserName: "test", LastName: "User"},
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "User not found",
			expectedCode:  http.StatusBadRequest,
			userID:        "2",
			updateDTO:     models.UpdateUserDTO{FirstName: "Updated", LastName: "User"},
			responseDTO:   nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				users.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(tc.responseDTO, tc.serviceError)
			}
			body, _ := json.Marshal(tc.updateDTO)
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodPut
			req.URL = httpSrv.URL + "/v1.0/users/" + tc.userID
			resp, err := req.SetBody(body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
