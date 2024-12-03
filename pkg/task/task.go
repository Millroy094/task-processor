package task

import (
	"encoding/json"
	"fmt"
	"time"
)

// EmailPayload represents the payload for an email task.
// @Description Payload for sending an email.
// @example {"email": "user@example.com", "subject": "Welcome!", "body": "Thank you for joining us!"}
type EmailPayload struct {
	Email   string `json:"email"`   // @example "user@example.com"
	Subject string `json:"subject"` // @example "Welcome to the Platform"
	Body    string `json:"body"`    // @example "Hello User, thank you for joining us."
}

// HealthCheckPayload represents the payload for performing a health check.
// @Description Payload for performing a health check.
// @example {"url": "https://example.com/health", "method": "GET", "status": "200 OK"}
type HealthCheckPayload struct {
	URL    string `json:"url"`    // @example "https://example.com/health"
	Method string `json:"method"` // @example "GET"
	Status string `json:"status"` // @example "200 OK"
}

// TaskResult represents the result of the task (e.g., status, error, detail).
type TaskResult struct {
	Status    string    `json:"status"`           // @example "completed"
	Error     string    `json:"error,omitempty"`  // @example "error occurred while sending email"
	Detail    string    `json:"detail,omitempty"` // @example "The email was sent successfully."
	Timestamp time.Time `json:"timestamp"`        // @example "2024-12-03T12:30:00Z"
}

// Task represents a task that can have different payloads based on its type.
// @Description Represents a task with dynamic payloads based on the task type.
// @Discriminator type
type Task struct {
	// Task ID
	// @example 1
	// @ReadOnly
	ID int `json:"id,omitempty" swaggerignore:"true"`

	// Task type determines the payload structure.
	// @example "email"
	Type string `json:"type" example:"email"`

	// Dynamic payload that varies by type.
	// @oneOf {EmailPayload} {HealthCheckPayload}
	// Use oneOf to indicate possible types of payload
	Payload interface{} `json:"payload"`

	// @example "2024-11-27T12:30:00Z"
	// @ReadOnly
	CreatedAt time.Time `json:"createdAt,omitempty" swaggerignore:"true"`

	// @example "pending"
	// @ReadOnly
	Status string `json:"status,omitempty" swaggerignore:"true"`

	// @example "2024-11-27T12:35:00Z"
	// @ReadOnly
	FinishedAt time.Time `json:"finishedAt,omitempty" swaggerignore:"true"`

	// @example {"status": "completed", "timestamp": "2024-12-03T12:30:00Z"}
	// @ReadOnly
	Result TaskResult `json:"result,omitempty" swaggerignore:"true"`
}

func (t *Task) UnmarshalPayload() (interface{}, error) {
	switch t.Type {
	case "email":
		var emailPayload EmailPayload
		if err := json.Unmarshal(t.Payload.([]byte), &emailPayload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal email payload: %w", err)
		}
		return emailPayload, nil

	case "health_check":
		var healthCheckPayload HealthCheckPayload
		if err := json.Unmarshal(t.Payload.([]byte), &healthCheckPayload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal health check payload: %w", err)
		}
		return healthCheckPayload, nil

	default:
		return nil, fmt.Errorf("unknown task type: %s", t.Type)
	}
}
