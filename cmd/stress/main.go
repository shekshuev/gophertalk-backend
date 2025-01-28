package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	baseURL := "http://localhost:3000/v1.0"
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
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Registration failed with status code: %d", resp.StatusCode)
	}
	resp.Body.Close()

	loginBody, _ := json.Marshal(map[string]string{
		"user_name": userName,
		"password":  password,
	})
	resp, err = http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		log.Fatalf("Failed to login user: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Login failed with status code: %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		log.Fatalf("Failed to decode login response: %v", err)
	}
	resp.Body.Close()

	accessToken := tokenResp.AccessToken
	if accessToken == "" {
		log.Fatal("Access token is empty")
	}

	maxPostsToView := 1000
	postsPerPage := 10
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := time.Second * 10

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

	targeter := vegeta.NewStaticTargeter(targets...)

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Feed Scrolling") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Total requests: %d\n", metrics.Requests)
	fmt.Printf("Average latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Max latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Min latency: %s\n", metrics.Latencies.Min)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Status codes: %v\n", metrics.StatusCodes)
}
