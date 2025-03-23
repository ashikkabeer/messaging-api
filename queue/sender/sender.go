package sender

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ashikkabeer/messaging-api/models"
	"github.com/rabbitmq/amqp091-go"
)

// golbal variable to hold the single sender instance
var (
	senderInstance *Sender
	once sync.Once
)

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
	var err error
	once.Do(func() {
		ch, chErr := conn.Channel()
		if chErr != nil {
			err = fmt.Errorf("failed to open channel: %v", chErr)
		return
	}

	// Declare a queue
	q, qerr := ch.QueueDeclare(
		"Messages",  
		false,       
		false,       
		false,       
		false,       
		nil,         
	)
	if qerr != nil {
		ch.Close()
		conn.Close()
		err = fmt.Errorf("failed to declare queue: %v", qerr)
			return

	}

	senderInstance = &Sender{
		Channel: ch,
		Queue:   q,
	}
	})
	return senderInstance, err

	
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

