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

func TestHandler_GetAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		limit         string
		offset        string
		responseDTOs  []models.ReadPostDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success fetching posts",
			expectedCode:  http.StatusOK,
			limit:         "10",
			offset:        "0",
			responseDTOs:  []models.ReadPostDTO{{ID: 1, Text: "Lorem Ipsum"}},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid limit",
			expectedCode:  http.StatusOK,
			limit:         "-1",
			offset:        "0",
			responseDTOs:  []models.ReadPostDTO{{ID: 1, Text: "Lorem Ipsum"}},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid offset",
			expectedCode:  http.StatusOK,
			limit:         "10",
			offset:        "-1",
			responseDTOs:  []models.ReadPostDTO{{ID: 1, Text: "Lorem Ipsum"}},
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
				posts.EXPECT().GetAllPosts(gomock.Any(), gomock.Any()).Return(tc.responseDTOs, tc.serviceError)
			}

			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + "/v1.0/posts?limit=" + tc.limit + "&offset=" + tc.offset
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_GetPostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		postID        string
		responseDTO   *models.ReadPostDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success fetching post",
			expectedCode:  http.StatusOK,
			postID:        "1",
			responseDTO:   &models.ReadPostDTO{ID: 1, Text: "Lorem Ipsum"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid post ID",
			expectedCode:  http.StatusNotFound,
			postID:        "abc",
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "Post not found",
			expectedCode:  http.StatusNotFound,
			postID:        "2",
			responseDTO:   nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().GetPostByID(gomock.Any()).Return(tc.responseDTO, tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + "/v1.0/posts/" + tc.postID
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		createDTO     models.CreatePostDTO
		responseDTO   *models.ReadPostDTO
		serviceError  error
		serviceCalled bool
	}{
		{
			name:         "Success create post",
			expectedCode: http.StatusCreated,
			createDTO: models.CreatePostDTO{
				Text: "Lorem ipsum",
			},
			responseDTO:   &models.ReadPostDTO{ID: 1, Text: "Lorem ipsum"},
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Validation error",
			expectedCode:  http.StatusUnprocessableEntity,
			createDTO:     models.CreatePostDTO{},
			responseDTO:   nil,
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:         "Error in service",
			expectedCode: http.StatusBadRequest,
			createDTO: models.CreatePostDTO{
				Text: "Lorem ipsum",
			},
			responseDTO:   nil,
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().CreatePost(gomock.Any()).Return(tc.responseDTO, tc.serviceError)
			}
			body, _ := json.Marshal(tc.createDTO)
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodPost
			req.URL = httpSrv.URL + "/v1.0/posts"
			resp, err := req.SetBody(body).Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_DeletePostByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		postID        string
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success deleting post",
			expectedCode:  http.StatusNoContent,
			postID:        "1",
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid post ID",
			expectedCode:  http.StatusNotFound,
			postID:        "abc",
			serviceError:  nil,
			serviceCalled: false,
		},
		{
			name:          "Post belongs to other user",
			expectedCode:  http.StatusNotFound,
			postID:        "2",
			serviceError:  assert.AnError,
			serviceCalled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().DeletePost(gomock.Any(), uint64(1)).Return(tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodDelete
			req.URL = httpSrv.URL + "/v1.0/posts/" + tc.postID
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_ViewPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		postID        string
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success view post",
			expectedCode:  http.StatusCreated,
			postID:        "1",
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid post ID",
			expectedCode:  http.StatusNotFound,
			postID:        "abc",
			serviceError:  nil,
			serviceCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().ViewPost(gomock.Any(), uint64(1)).Return(tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodPost
			req.URL = httpSrv.URL + "/v1.0/posts/" + tc.postID + "/view"
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_LikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		postID        string
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success like post",
			expectedCode:  http.StatusCreated,
			postID:        "1",
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid post ID",
			expectedCode:  http.StatusNotFound,
			postID:        "abc",
			serviceError:  nil,
			serviceCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().LikePost(gomock.Any(), uint64(1)).Return(tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodPost
			req.URL = httpSrv.URL + "/v1.0/posts/" + tc.postID + "/like"
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}

func TestHandler_DislikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	posts := mocks.NewMockPostService(ctrl)
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
	handler := NewHandler(nil, nil, posts, &cfg)
	httpSrv := httptest.NewServer(handler.Router)
	defer httpSrv.Close()

	testCases := []struct {
		name          string
		expectedCode  int
		postID        string
		serviceError  error
		serviceCalled bool
	}{
		{
			name:          "Success dislike post",
			expectedCode:  http.StatusNoContent,
			postID:        "1",
			serviceError:  nil,
			serviceCalled: true,
		},
		{
			name:          "Invalid post ID",
			expectedCode:  http.StatusNotFound,
			postID:        "abc",
			serviceError:  nil,
			serviceCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.serviceCalled {
				posts.EXPECT().DislikePost(gomock.Any(), uint64(1)).Return(tc.serviceError)
			}
			req := resty.New().R()
			req.Header.Set("Authorization", "Bearer "+accessToken)
			req.Method = http.MethodDelete
			req.URL = httpSrv.URL + "/v1.0/posts/" + tc.postID + "/like"
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
