package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	srvs "portqueue/services"
)

func init() {
	err := srvs.LoadCompletedRequests()

	if err != nil {
		fmt.Println("Error loading completed requests:", err)
	}
}

func main() {
	requests, err := srvs.GetRequestData("requests.json")

	if err != nil {
		fmt.Println("Error reading requests.json:", err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	for _, request := range requests {
		if !srvs.CompletedRequests[request.Id] {
			err := srvs.ExecuteRequest(request)
			if err != nil {
				fmt.Println("Error executing request:", err)
			}
		}
	}

	<-signalChan // Wait for SIGINT

	srvs.SaveCompletedRequests()
}

