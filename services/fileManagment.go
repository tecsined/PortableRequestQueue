package services

import (
	"encoding/json"
	"fmt"
	"os"
)

func init() {
	if err := LoadCompletedRequests(); err != nil {
		fmt.Println("Error loading completed requests:", err)
	}
}

func GetRequestData(filePath string) ([]RequestData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var requests []RequestData
	err = json.Unmarshal(data, &requests)
	if err != nil {
		return nil, err
	}

	return requests, err
}

const completedRequestsFile = "completed_requests.json"

var CompletedRequests = map[string]bool{}

func LoadCompletedRequests() error {
	data, err := os.ReadFile(completedRequestsFile)
	if err != nil {
		return nil // If file doesn't exist, start with an empty map
	}

	err = json.Unmarshal(data, &CompletedRequests)
	if err != nil {
		return err
	}
	return nil
}

func SaveCompletedRequests() {
	data, err := json.Marshal(CompletedRequests)
	if err != nil {
		fmt.Println("Error saving completed requests:", err)
		return
	}
	const FILE_PERMISSION = 0644
	err = os.WriteFile(completedRequestsFile, data, FILE_PERMISSION)
	if err != nil {
		fmt.Println("Error saving completed requests:", err)
	}

	fmt.Println("Completed requets saved")
}
