package task

import (
	"encoding/json"
	"time"
)

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

// @Description Represents a task with ID, type, and payload information.
type Task struct {
	// @ID Task ID
	// @example 1
	// It will only appear when retrieving the task (GET or PUT).
	// @ReadOnly
	ID int `json:"id"`

	// @Type Type of the task (e.g., "email", "health check")
	// @example "email"
	Type string `json:"type"`

	// @Payload The actual task payload (which can vary depending on the task type)
	// Using RawMessage to delay unmarshalling
	Payload json.RawMessage `json:"payload"`

	// @CreatedAt The time the task was created
	// @example "2024-11-27T12:30:00Z"
	// It will only appear when retrieving the task (GET or PUT).
	// @ReadOnly
	CreatedAt time.Time `json:"createdAt"`

	// @Status The current status of the task (e.g., "pending", "in-progress", "completed")
	// @example "pending"
	// This field will be auto-populated on the server-side, hence not included in the body when creating a task.
	// It will only appear when retrieving the task (GET or PUT).
	// @ReadOnly
	Status string `json:"status"`

	// @FinishedAt The time when the task was finished (only updated when task is completed)
	// @example "2024-11-27T12:35:00Z"
	// This field will be auto-populated on the server-side, hence not included in the body when creating a task.
	// It will only appear when retrieving the task (GET or PUT).
	// @ReadOnly
	FinishedAt time.Time `json:"finishedAt,omitempty"`

	// @Result The result of the task, includes status, error, and detail.
	// This field will be auto-populated by the backend after task completion.
	// @example {"status": "completed", "timestamp": "2024-12-03T12:35:00Z"}
	// @ReadOnly
	Result TaskResult `json:"result,omitempty"` // This field will be set by the server, not the client.
}
