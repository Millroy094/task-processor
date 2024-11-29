package task

// @Description Represents a task with ID, type, and payload information.
type Task struct {
	// @ID Task ID
	// @example 1
	ID int `json:"id"`

	// @Type Type of the task (e.g., "email", "health check")
	// @example "email"
	Type string `json:"type"`

	// @Payload The actual task payload (which can vary depending on the task type)
	// @example "{\"email\":\"recipient@example.com\",\"subject\":\"Welcome!\",\"body\":\"Hello!\"}"
	Payload string `json:"payload"`
}
