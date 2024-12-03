package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/millroy094/task-processor/pkg/task"
)

type HealthCheckPayload struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Status string `json:"status"`
}

func performHealthCheck(task task.Task) error {
	var data HealthCheckPayload
	err := json.Unmarshal([]byte(task.Payload), &data)
	if err != nil {
		log.Printf("Failed to unmarshal health check payload: %v", err)
		return err
	}

	expectedStatusCode, err := strconv.Atoi(data.Status)
	if err != nil {
		return fmt.Errorf("invalid status code: %v", err)
	}

	req, err := http.NewRequest(data.Method, data.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform health check: %v", err)
	}
	defer resp.Body.Close()

	result := task.Result
	result.Status = "completed"
	result.Timestamp = time.Now()

	if resp.StatusCode != expectedStatusCode {
		result.Status = "failed"
		result.Error = fmt.Sprintf("expected status %d, got %d", expectedStatusCode, resp.StatusCode)
	} else {
		result.Detail = fmt.Sprintf("Health check passed for %s: received expected status %d", data.URL, resp.StatusCode)
	}

	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"result":    result,
			"updatedAt": time.Now(),
		},
	}
	_, err = taskCollection.UpdateOne(context.Background(), map[string]interface{}{"id": task.ID}, update)
	if err != nil {
		log.Printf("Error updating task result: %v", err)
		return err
	}

	return nil
}
