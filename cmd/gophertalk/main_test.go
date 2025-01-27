package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestFeedScrollingWithVegeta(t *testing.T) {
	baseURL := "http://localhost:3000/v1.0"

	// Регистрация и авторизация пользователя
	userName := fmt.Sprintf("testuser_%d", time.Now().Unix())
	password := "TestPassword123!"
	registerBody, _ := json.Marshal(map[string]string{
		"user_name":        userName,
		"password":         password,
		"password_confirm": password,
		"first_name":       "Test",
		"last_name":        "User",
	})
	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(registerBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Авторизация для получения токена
	loginBody, _ := json.Marshal(map[string]string{
		"user_name": userName,
		"password":  password,
	})
	resp, err = http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	assert.NoError(t, err)
	resp.Body.Close()

	// Токен для авторизации
	accessToken := tokenResp.AccessToken
	assert.NotEmpty(t, accessToken)

	// Конфигурация теста
	maxPostsToView := 1000                           // Лимит на количество постов для просмотра
	postsPerPage := 10                               // Количество постов на одной странице
	rate := vegeta.Rate{Freq: 100, Per: time.Second} // Частота запросов (10 в секунду)
	duration := time.Second * 10                     // Продолжительность теста

	// Создаем цели для GET-запросов
	var targets []vegeta.Target
	for offset := 0; offset < maxPostsToView; offset += postsPerPage {
		targets = append(targets, vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/posts?limit=%d&offset=%d", baseURL, postsPerPage, offset),
			Header: map[string][]string{
				"Authorization": {"Bearer " + accessToken},
			},
		})
	}

	// Создаем статический таргетер
	targeter := vegeta.NewStaticTargeter(targets...)

	// Запуск Vegeta для замера GET-запросов
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Feed Scrolling") {
		metrics.Add(res)
	}
	metrics.Close()

	// Отчет по GET-запросам
	fmt.Printf("Total requests: %d\n", metrics.Requests)
	fmt.Printf("Average latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Max latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Min latency: %s\n", metrics.Latencies.Min)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Status codes: %v\n", metrics.StatusCodes)
}
