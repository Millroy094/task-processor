package task

import (
	"encoding/json"
	"time"
)

// @Description Represents a task with ID, type, and payload information.
type Task struct {
	// @ID Task ID
	// @example 1
	ID int `json:"id"`

	// @Type Type of the task (e.g., "email", "health check")
	// @example "email"
	Type string `json:"type"`

	// @Payload The actual task payload (which can vary depending on the task type)
	// Using RawMessage to delay unmarshalling
	Payload json.RawMessage `json:"payload"`

	// @CreatedAt The time the task was created
	// @example "2024-11-27T12:30:00Z"
	CreatedAt time.Time `json:"createdAt"`

	// @Status The current status of the task (e.g., "pending", "in-progress", "completed")
	// @example "pending"
	Status string `json:"status"`

	// @FinishedAt The time when the task was finished (only updated when task is completed)
	// @example "2024-11-27T12:35:00Z"
	FinishedAt time.Time `json:"finishedAt,omitempty"`
}
