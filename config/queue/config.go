package queue

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ashikkabeer/messaging-api/queue/receiver"
	"github.com/ashikkabeer/messaging-api/queue/sender"
	"github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string 
	Port     int   
	User     string 
	Password string
}

type Queue struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel   
	Queue   amqp091.Queue
}

// Connection wraps the RabbitMQ connection
type Connection struct {
	Conn *amqp091.Connection 
}

// Global variables to maintain singleton instances of queue connections and handlers
var (
	queueConnection *Connection      
	queueSender    *sender.Sender   
	queueReceiver  *receiver.Receiver 
)


// Sets up the RabbitMQ connection
// creates sender and receiver instances
func InitializeQueue() error {
	var err error
	// Create the initial RabbitMQ connection
	queueConnection, err = CreateConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Initialized sender
	queueSender, err = sender.NewSender(queueConnection.Conn)
	if err != nil {
		queueConnection.Conn.Close()
		return fmt.Errorf("failed to create sender: %v", err)
	}
	sender.SetSenderInstance(queueSender)

	// Initialized receiver
	queueReceiver, err = receiver.NewReceiver(queueConnection.Conn)
	if err != nil {
		queueSender.Close()
		queueConnection.Conn.Close()
		return fmt.Errorf("failed to create receiver: %v", err)
	}

	// Start consuming messages from the queue
	if err := queueReceiver.StartConsuming(); err != nil {
		queueReceiver.Close()
		queueSender.Close()
		queueConnection.Conn.Close()
		return fmt.Errorf("failed to start consuming: %v", err)
	}

	return nil
}


func CreateConnection() (*Connection, error) {
	config := NewConfig()
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	conn, err := amqp091.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitmq connection: %v", err)
	}
	return &Connection{Conn: conn}, nil
}

//  closes all queue related connections
func CloseConnections() {
	if queueReceiver != nil {
		queueReceiver.Close()
	}
	if queueSender != nil {
		queueSender.Close()
	}
	if queueConnection != nil && queueConnection.Conn != nil {
		queueConnection.Conn.Close()
	}
}

func NewConfig() *Config {
	port, _ := strconv.Atoi(getEnvOrDefault("RabbitMQ_PORT", "5672"))
	return &Config{
		Host:     getEnvOrDefault("RabbitMQ_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("RabbitMQ_USER", "guest"),
		Password: getEnvOrDefault("RabbitMQ_PASSWORD", "guest"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

