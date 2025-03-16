package utils

import (
    "fmt"
    "time"
    "github.com/rabbitmq/amqp091-go"
)

func WaitForRabbitMQ() {
    fmt.Println("Waiting for RabbitMQ to be ready...")
    for i := 0; i < 60; i++ {
        conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
        if err == nil {
            fmt.Println("RabbitMQ is ready!")
            conn.Close() 
            return      
        }
        fmt.Printf("RabbitMQ not ready yet, waiting... (attempt %d/12)\n", i+1)
        time.Sleep(1 * time.Second)
    }
    fmt.Println("Warning: Timeout waiting for RabbitMQ")
}