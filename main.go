package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GenerateRequest struct {
	SystemPrompt string `json:"system_prompt"`
	LLM          string `json:"LLM"`
	Text         string `json:"text"`
	TaskID       string `json:"task_id"`
}

type GenerateResponse struct {
	Status       bool   `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type AIRequest struct {
	SystemPrompt string `json:"system_prompt"`
	LLM          string `json:"LLM"`
	Text         string `json:"text"`
}

type AIResponse struct {
	Response string `json:"response"`
}

type CallbackRequest struct {
	TaskID   string `json:"task_id"`
	Response string `json:"response"`
}

var (
	callbackURL = os.Args[1]
	aiURL       = os.Args[2]
)

func main() {
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/callback-test", callbackTestHandler)

	log.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest

	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	log.Println("Received start request:", req.TaskID)

	// --- HEALTH CHECK ИИ-СЕРВИСА ---
	if !checkAIHealth() {
		log.Println("AI service is NOT available")

		resp := GenerateResponse{
			Status:       false,
			ErrorMessage: "AI service is not reachable",
		}
		out, _ := json.Marshal(resp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write(out)
		return
	}
	// --------------------------------

	resp := GenerateResponse{Status: true, ErrorMessage: ""}
	out, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

	req.Text = strings.ReplaceAll(req.Text, "\r\n", " ")
	req.Text = strings.ReplaceAll(req.Text, "\n", " ")

	go processAsync(req)
}

// Тестовый endpoint чтобы видеть callback
func callbackTestHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Println("Got callback:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": false}`))
}
