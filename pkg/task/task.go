package task

import (
	"encoding/json"
	"fmt"
	"time"
)

type EmailPayload struct {
	Email   string `json:"email"`   
	Subject string `json:"subject"`
	Body    string `json:"body"`   
}

type HealthCheckPayload struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Status string `json:"status"`
}

type TaskRequest struct {
	Type string `json:"type" example:"email"`
	Payload json.RawMessage `json:"payload"`
}

type TaskResult struct {
	Status    string    `json:"status"`      
	Error     string    `json:"error,omitempty"` 
	Detail    string    `json:"detail,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Task struct {
	ID int `json:"id,omitempty"`
	Type string `json:"type" example:"email"`
	Payload json.RawMessage `json:"payload"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Status string `json:"status,omitempty"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Result TaskResult `json:"result,omitempty"`
}

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

