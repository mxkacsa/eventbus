package main

import (
	"encoding/json"
	"fmt"
	"github.com/mxkacsa/eventbus"
	"io"
	"log"
	"net/http"
)

type Data struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

var whService *WebhookService

type WebhookService struct {
	Eb *eventbus.EventBus[Data]
}

func NewWebhookService() *WebhookService {
	return &WebhookService{
		Eb: new(eventbus.EventBus[Data]),
	}
}

func main() {
	whService = NewWebhookService()
	workerService := SampleWorkerService{
		WhService: whService,
	}

	err := workerService.Start()
	if err != nil {
		return
	}

	http.HandleFunc("/webhook", handlePost)
	log.Println("Server running on port 8080")
	panic(http.ListenAndServe(":8080", nil))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing json", http.StatusBadRequest)
		return
	}

	//whService.Eb.Publish(data)
	whService.Eb.PublishAsync(data)

	_, err = fmt.Fprintf(w, "Message processed")
	if err != nil {
		return
	}
}

// Client side

type SampleWorkerService struct {
	WhService *WebhookService
}

func (s *SampleWorkerService) Start() error {
	return s.WhService.Eb.SubscribeMethod(s, s.process)
}

func (s *SampleWorkerService) Stop() error {
	return s.WhService.Eb.UnsubscribeMethod(s, s.process)
}

func (s *SampleWorkerService) process(data Data) {
	log.Println("data received:", data.Name, data.Message)
}
