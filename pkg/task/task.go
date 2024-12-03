package task

import (
	"encoding/json"
	"fmt"
	"time"
)

// @Description Represents the payload for an email task.
type EmailPayload struct {
	// @example "user@example.com"
	Email string `json:"email"`

	// @example "Welcome to the Platform"
	Subject string `json:"subject"`

	// @example "Hello User, thank you for joining us."
	Body string `json:"body"`
}

// @Description Represents the payload for a health check task.
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

// @Description Represents a task with type and payload information.
// @DiscriminatorField type
// @DiscriminatorMapping email EmailPayload, health_check HealthCheckPayload
type Task struct {
	// @ID Task ID
	// @example 1
	// It will only appear when retrieving the task (GET).
	// @ReadOnly
	ID int `json:"id,omitempty" swaggerignore:"true"`

	// @Type Type of the task (e.g., "email", "health_check")
	// @example "email"
	Type string `json:"type"`

	// @Payload The actual task payload, which varies based on the task type.
	// Using RawMessage to delay unmarshalling.
	// @oneOf {object} EmailPayload {object} HealthCheckPayload
	Payload json.RawMessage `json:"payload"`

	// @CreatedAt The time the task was created.
	// @example "2024-11-27T12:30:00Z"
	// It will only appear when retrieving the task (GET).
	// @ReadOnly
	CreatedAt time.Time `json:"createdAt,omitempty" swaggerignore:"true"`

	// @Status The current status of the task (e.g., "pending", "in-progress", "completed").
	// @example "pending"
	// This field will be auto-populated on the server-side.
	// It will only appear when retrieving the task (GET).
	// @ReadOnly
	Status string `json:"status,omitempty" swaggerignore:"true"`

	// @FinishedAt The time when the task was finished (only updated when task is completed).
	// @example "2024-11-27T12:35:00Z"
	// This field will be auto-populated on the server-side.
	// It will only appear when retrieving the task (GET).
	// @ReadOnly
	FinishedAt time.Time `json:"finishedAt,omitempty" swaggerignore:"true"`

	// @Result The result of the task, includes status, error, and detail.
	// This field will be auto-populated by the backend after task completion.
	// @example {"status": "completed", "timestamp": "2024-12-03T12:35:00Z"}
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
