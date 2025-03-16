package sender

import (
	"encoding/json"
	"fmt"

	"github.com/ashikkabeer/messaging-api/models"
	"github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewSender() (*Sender, error) {
	// Connect to RabbitMQ
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare a queue
	q, err := ch.QueueDeclare(
		"Messages",  
		false,       
		false,       
		false,       
		false,       
		nil,         
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	return &Sender{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (s *Sender) SendMessage(message models.RequestBody) error {

	jsonData, err := json.Marshal(message)
	if err!= nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	err = s.channel.Publish(
		"",          
		"Messages", 
		false,        
		false,       
		amqp091.Publishing{
			ContentType: "application/json",
			Body:       	jsonData,

		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (s *Sender) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}

