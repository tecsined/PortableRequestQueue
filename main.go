package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	srvs "portqueue/services"
)

func init() {
	err := srvs.LoadCompletedRequests()

	if err != nil {
		fmt.Println("Error loading completed requests:", err)
	}
}

func main() {

	displaySummary()
	startProcess := processRequestQuestion()

	if !startProcess {
		return
	}

	requests, err := srvs.GetRequestData("requests.json")

	if err != nil {
		fmt.Println("Error reading requests.json:", err)
		return
	}

	for _, request := range requests {
		if !srvs.CompletedRequests[request.Id] {
			err := srvs.ExecuteRequest(request)
			if err != nil {
				fmt.Println("Error executing request:", err)
			}
		}
	}

	srvs.SaveCompletedRequests()
}

func displaySummary() {
	fmt.Println("Queue processor works using a request list stored in a json file name requests.json\n")
	fmt.Println("Make sure that the file exist in the root of the application\n")
	fmt.Println("***********************************************************************************\n")
}

func processRequestQuestion() bool {
	const yes = "y"
	const no = "n"
	fmt.Println("Do you want to start the processing the request. Please type y for yes\n")
	fmt.Println("If you dont type y the program will end\n")
	in := bufio.NewReader(os.Stdin)
	choice, _ := in.ReadString('\n')
	choice = strings.ToLower(strings.TrimSpace(choice))
	return choice == yes
}
