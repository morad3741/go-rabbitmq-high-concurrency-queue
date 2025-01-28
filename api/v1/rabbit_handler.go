package v1

import (
	"encoding/json"
	"fmt"
	"hexArchitectureProject/internal/mq"
	"net/http"
	"time"
)

// SendMessageRequest represents the expected request structure for sending a message.
type SendMessageRequest struct {
	QueueName string `json:"queueName"`
	Event     string `json:"event"`
}

// RabbitMQHandler handles HTTP requests related to RabbitMQ operations.
type RabbitMQHandler struct {
	rabbitMQService *mq.RabbitMQService
}

// NewRabbitMQHandler creates a new instance of RabbitMQHandler with dependency injection.
func NewRabbitMQHandler(rabbitMQService *mq.RabbitMQService) *RabbitMQHandler {
	return &RabbitMQHandler{
		rabbitMQService: rabbitMQService,
	}
}

// SendMessageToRabbit handles the HTTP request for sending a message to RabbitMQ.
func (h *RabbitMQHandler) SendMessageToRabbit(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest

	// Parse JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.QueueName == "" || req.Event == "" {
		http.Error(w, "QueueName and Event are required", http.StatusBadRequest)
		return
	}

	// Ensure queue is defined
	// Adjust numProducerThreads and numConsumerThreads as needed, or pass dynamic values if necessary
	err := h.rabbitMQService.DefineQueue(req.QueueName, 3, 3, func(event string) {
		// Print message and start sleep
		fmt.Printf("Received message: %s\nStarting sleep for 10 seconds...\n", event)

		// Sleep for 10 seconds
		time.Sleep(10 * time.Second)

		// Print message after sleep
		fmt.Printf("Finished sleep for 10 seconds, Processed message: %s\n", event)
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to define queue: %v", err), http.StatusInternalServerError)
		return
	}

	// Send message to RabbitMQ
	err = h.rabbitMQService.SendMessage(req.QueueName, req.Event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond to client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent successfully"))
}
