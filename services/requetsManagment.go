package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func init() {
	err := LoadCompletedRequests()

	if err != nil {
		fmt.Println("Error loading completed requests:", err)
	}
}

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
			trackRetries(request, attempt, maxRetries, retryDelay)
			continue
		}
	}
	return nil
}

func trackRetries(request RequestData, attempt int, maxRetries int, retryDelay time.Duration) {
	fmt.Printf("Retrying request to %s (attempt %d)...\n", request.URL, attempt+1)
	if attempt == maxRetries-1 {
		fmt.Printf("Retrying ended for this URL and process will continue if there are more requests")
	}
	time.Sleep(retryDelay)
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

func NewRequestsWorker() <-chan RequestData {
	var outCh = make(chan RequestData)
	requests, err := GetRequestData("requests.json")
	if err != nil {
		panic(fmt.Errorf("error reading requests.json: %w", err))
	}

	go func() {
		for _, request := range requests {
			outCh <- request
		}
		close(outCh)
	}()
	return outCh
}

const DEFAULT_CONCURRENCY = 5

func StartCrawling() {
	var wg sync.WaitGroup
	newRequestsCh := NewRequestsWorker()

	for range DEFAULT_CONCURRENCY {
		wg.Add(1)
		go func(reqsCh <-chan RequestData) {
			for r := range reqsCh {
				if !CompletedRequests[r.Id] {
					err := ExecuteRequest(r)
					if err != nil {
						fmt.Println("Error executing request:", err)
					}
				}
			}
			wg.Done()
		}(newRequestsCh)
		wg.Wait()
	}
}
