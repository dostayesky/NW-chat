package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/wutthichod/sa-connext/shared/config"
	"github.com/wutthichod/sa-connext/shared/messaging"
)

func main() {
	godotenv.Load("../.env") // ./ = โฟลเดอร์เดียวกับ main.go
	config := config.LoadConfig()
	rb, err := messaging.NewRabbitMQ(config.RabbitURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rb.Close()

	queueName, err := rb.SetupQueue(
		"email_queue",              // queue name
		"notification.exchange",    // exchange
		"direct",                   // exchange type
		"notification.email",       // routing key
		true,                       // durable
		nil,                        // args
	)
	if err != nil {
    	log.Fatalf("Failed to setup email queue: %v", err)
	}

	emailConsumer := messaging.NewEmailConsumer(rb, queueName, config.Email, config.EmailPW)
	if err := emailConsumer.Start(); err != nil {
		log.Fatalf("Failed to start email consumer: %v", err)
	}
	select {} // Block forever
}