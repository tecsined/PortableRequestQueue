package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestData struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body,omitempty"`
}

type Requests []RequestData

const DEFAULT_RETRIES = 5

func ExecuteRequest(request RequestData) error {
	maxRetries := DEFAULT_RETRIES
	retryDelay := 1 * time.Second

	req, err := BuildHttpRequest(request)
	if err != nil {
		return err
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := SendHttpRequest(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			CompletedRequests[request.Id] = true
			SaveCompletedRequests()
			err = HandleResponseBody(resp)
			return err
		} else {
			fmt.Printf("Retrying request to %s (attempt %d)...\n", request.URL, attempt+1)
			if attempt == maxRetries-1 {
				fmt.Printf("Retrying ended for this URL and process will continue if there are more requests")
			}
			time.Sleep(retryDelay)
			continue
		}
	}
	return nil
}

func BuildHttpRequest(request RequestData) (*http.Request, error) {
	body := []byte(request.Body)
	return http.NewRequest(request.Method, request.URL, bytes.NewReader(body))
}

func SendHttpRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func HandleResponseBody(resp *http.Response) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Println("Method:", resp.Request.Method)
	fmt.Println("URL:", resp.Request.URL)
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(responseBody))
	fmt.Println("----------")
	return nil
}
