package main

import (
	"net/http"
	"time"
)

func checkAIHealth() bool {
	healthURL := aiURL // ← вставь свой health endpoint

	client := http.Client{
		Timeout: 3 * time.Second, // быстрый timeout
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}
