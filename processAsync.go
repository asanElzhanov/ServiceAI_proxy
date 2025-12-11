package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// Callback
func sendCallback(taskId string, response string) {
	cb := CallbackRequest{
		TaskID:   taskId,
		Response: response,
	}

	jsonCb, _ := json.Marshal(cb)

	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(jsonCb))
	if err != nil {
		log.Println("Callback failed:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Callback sent. Status:", resp.StatusCode)
}

func processAsync(req GenerateRequest) {
	log.Println("Processing task:", req.TaskID)


	aiReq := AIRequest{
		SystemPrompt: req.SystemPrompt,
		LLM:          req.LLM,
		Text:         req.Text,
	}

	jsonReq, _ := json.Marshal(aiReq)

	httpClient := &http.Client{
		Timeout: 30 * time.Minute, 
	}

	aiRespRaw, err := httpClient.Post(
		aiURL, 
		"application/json",
		bytes.NewBuffer(jsonReq),
	)
	if err != nil {
		log.Println("AI request error:", err)
		sendCallback(req.TaskID, "ERROR: "+err.Error())
		return
	}
	defer aiRespRaw.Body.Close()

	var aiResp AIResponse
	body, _ := io.ReadAll(aiRespRaw.Body)
	json.Unmarshal(body, &aiResp)

	log.Println("AI response:", aiResp.Response)

	// Отправить callback
	sendCallback(req.TaskID, aiResp.Response)
}