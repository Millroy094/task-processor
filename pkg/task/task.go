package task

// @Description Represents a task with ID, type, and payload information.
type Task struct {
    // @ID Task ID
    // @example 1
    ID      int    `json:"id"`

    // @Type Type of the task
    // @example "email"
    Type    string `json:"type"`

    // @Payload The actual task payload
    // @example "Send welcome email"
    Payload string `json:"payload"`
}