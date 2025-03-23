package main

import (
	"log"

	"github.com/ashikkabeer/messaging-api/api/routes"
	"github.com/ashikkabeer/messaging-api/config/db"
	"github.com/ashikkabeer/messaging-api/config/queue"
	"github.com/ashikkabeer/messaging-api/utils"
)

func main() {
	// waiting for rabbitmq to start
	utils.WaitForRabbitMQ()

	// setup router and db
	r := routes.SetupRouter()
	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}
	if err:= queue.InitializeQueue(); err != nil {
		log.Fatalf("failed to initialize queue",err)
	}

	defer queue.CloseConnections()
	r.Run()
}