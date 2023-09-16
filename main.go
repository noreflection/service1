// Microservice 1: main.go

package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

// Message Define a struct to represent the message
type Message struct {
	Text string `json:"text"`
}

// RabbitMQ connection parameters
const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "messages"
)

func main() {
	http.HandleFunc("/receive", receiveMessageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func receiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "Failed to decode message", http.StatusBadRequest)
		return
	}

	// Publish the message to RabbitMQ
	err := publishToRabbitMQ(message)
	if err != nil {
		http.Error(w, "Failed to send message to RabbitMQ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Message received and sent to RabbitMQ"))
}

func publishToRabbitMQ(message Message) error {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("cannot connect")
		return err
	}
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(ch)

	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,     // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
