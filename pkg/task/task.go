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
	// @example "user@example.com"
	Email string `json:"email"`

	// @example "Welcome to the Platform"
	Subject string `json:"subject"`

	// @example "Hello User, thank you for joining us."
	Body string `json:"body"`
}

// @Description Payload for performing a health check.
// @example {"url": "https://example.com/health", "method": "GET", "status": "200 OK"}
type HealthCheckPayload struct {
	// @example "https://example.com/health"
	URL string `json:"url"`

	// @example "GET"
	Method string `json:"method"`

	// @example "200 OK"
	Status string `json:"status"`
}

// TaskResult represents the result of the task (e.g., status, error, detail).
type TaskResult struct {
	// @example "completed"
	Status string `json:"status"`

	// @example "error occurred while sending email"
	Error string `json:"error,omitempty"`

	// @example "The email was sent successfully."
	Detail string `json:"detail,omitempty"`

	// @example "2024-12-03T12:30:00Z"
	Timestamp time.Time `json:"timestamp"`
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
	// @oneOf {object} EmailPayload {object} HealthCheckPayload
	// @example {"email": "user@example.com", "subject": "Welcome!", "body": "Thank you for joining us!"}
	Payload json.RawMessage `json:"payload"`

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

// UnmarshalPayload deserializes the task payload based on its type.
func (t *Task) UnmarshalPayload() (interface{}, error) {
	switch t.Type {
	case "email":
		var emailPayload EmailPayload
		if err := json.Unmarshal(t.Payload, &emailPayload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal email payload: %w", err)
		}
		return emailPayload, nil

	case "health_check":
		var healthCheckPayload HealthCheckPayload
		if err := json.Unmarshal(t.Payload, &healthCheckPayload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal health check payload: %w", err)
		}
		return healthCheckPayload, nil

	default:
		return nil, fmt.Errorf("unknown task type: %s", t.Type)
	}
}
