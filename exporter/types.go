package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JobItem represents a single job definition from the Jenkins jobs API.
type JobItem struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	URL   string `json:"url"`
}

// QueueItemTask represents a single task definition in a queue
type QueueItemTask struct {
	Name string `json:"name"`
}

// QueueItem represents a queue item definition from the Jenkins queue API.
type QueueItem struct {
	Task         *QueueItemTask `json:"task"`
	InQueueSince float64        `json:"inQueueSince"`
}

// Collector represents the collected api response from the Jenkins APIs.
type Collector struct {
	Jobs  []*JobItem   `json:"jobs"`
	Queue []*QueueItem `json:"items"`
}

// Fetch gathers the content from the Jenkins API.
func (c *Collector) Fetch(address, username, password string) error {
	reqJobs, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/json", address),
		nil,
	)

	if username != "" && password != "" {
		reqJobs.SetBasicAuth(username, password)
	}

	jobs, err := simpleClient().Do(reqJobs)
	if err != nil {
		return fmt.Errorf("failed to request jobs api. %s", err)
	}
	defer jobs.Body.Close()

	if err := json.NewDecoder(jobs.Body).Decode(c); err != nil {
		return fmt.Errorf("failed to parse jobs api. %s", err)
	}

	reqQueue, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/queue/api/json", address),
		nil,
	)

	if username != "" && password != "" {
		reqQueue.SetBasicAuth(username, password)
	}

	queue, err := simpleClient().Do(reqQueue)
	if err != nil {
		return fmt.Errorf("failed to request queue api. %s", err)
	}
	defer queue.Body.Close()

	if err := json.NewDecoder(queue.Body).Decode(c); err != nil {
		return fmt.Errorf("failed to parse queue api. %s", err)
	}

	return nil
}
