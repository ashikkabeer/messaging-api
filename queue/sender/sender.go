package sender

import (
	"encoding/json"
	"fmt"

	"github.com/ashikkabeer/messaging-api/models"
	"github.com/rabbitmq/amqp091-go"
)

// golbal variable to hold the single sender instance
var senderInstance *Sender

func SetSenderInstance(s *Sender) {
    senderInstance = s
}
type Sender struct {
	// Conn    *amqp091.Connection
	Channel *amqp091.Channel
	Queue   amqp091.Queue
}

func NewSender(conn *amqp091.Connection) (*Sender, error) {
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
		Channel: ch,
		Queue:   q,
	}, nil
}
func SendMessageToQueue(message models.RequestBody) error {
    // Check if sender is initialized
    if senderInstance == nil {
        return fmt.Errorf("default sender not initialized")
    }
    // Use the singleton instance to send message
    return senderInstance.SendMessage(message)
}


func (s *Sender) SendMessage(message models.RequestBody) error {
	jsonData, err := json.Marshal(message)
	if err!= nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	err = s.Channel.Publish(
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
	if s.Channel != nil {
		s.Channel.Close()
	}
}

