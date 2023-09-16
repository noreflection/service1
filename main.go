// microservice1.go

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Define a global variable to store the message temporarily
var message string

// Define an endpoint to receive messages via POST requests
func receiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	message = string(body)
	fmt.Fprintf(w, "Message received: %s", message)

	// Pass the message to Microservice 2
	passMessageToMicroservice2(message)
}

// Function to pass the message to Microservice 2
func passMessageToMicroservice2(message string) {
	// Send an HTTP POST request to Microservice 2 with the message
	_, err := http.Post("http://microservice2:8080/receive", "application/json", bytes.NewBuffer([]byte(message)))
	if err != nil {
		log.Printf("Failed to pass message to Microservice 2: %v", err)
	}
}

func main() {
	http.HandleFunc("/receive", receiveMessageHandler)
	//fmt.Print("hey")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
