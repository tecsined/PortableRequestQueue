package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	srcs "portqueue/services"
)

func main() {

	displayInstructions()

	if !shouldStartCrawling() {
		return
	}

	srcs.StartCrawling()

	srcs.SaveCompletedRequests()
}

func displayInstructions() {
	fmt.Println("Queue processor works using a request list stored in a json file name requests.json")
	fmt.Println("Make sure that the file exist in the root of the application")
	fmt.Println("***********************************************************************************")
}

func shouldStartCrawling() bool {
	const yes = "y"
	fmt.Println("Do you want to start the processing the request. Please type y for yes")
	fmt.Println("If you dont type y the program will end")
	in := bufio.NewReader(os.Stdin)
	choice, _ := in.ReadString('\n')
	choice = strings.ToLower(strings.TrimSpace(choice))
	return choice == yes
}
