package main

import (
	"github.com/ashikkabeer/messaging-api/api/routes"
	"github.com/ashikkabeer/messaging-api/config/db"
	"github.com/ashikkabeer/messaging-api/queue/receiver"
	"log"
	"github.com/ashikkabeer/messaging-api/utils"
)

func main() {
	// waiting for rabbitmq to start
	utils.WaitForRabbitMQ()

	// setup router and db
	r := routes.SetupRouter()
	db.Connect()

	// start consuming messages from rabbitmq
	messageReceiver, err := receiver.NewReceiver()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer messageReceiver.Close()
	if err := messageReceiver.StartConsuming();err != nil {
        log.Fatalf("Failed to start consuming: %v", err)
    }
	// if err := messageReceiver.StartConsuming(receiver.ProcessMessage); err != nil {
    //     log.Fatalf("Failed to start consuming: %v", err)
    // }

	r.Run()
}
