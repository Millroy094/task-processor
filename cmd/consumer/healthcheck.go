package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf("health check failed for %s: expected status %d, got %d", data.URL, expectedStatusCode, resp.StatusCode)
	}

	fmt.Printf("Health check passed for %s: received expected status %d\n", data.URL, resp.StatusCode)
	return nil
}
