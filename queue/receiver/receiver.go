package receiver

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/ashikkabeer/messaging-api/config/db"
	"github.com/ashikkabeer/messaging-api/models"
	"fmt"
	"encoding/json"
    "log"
)

type Receiver struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewReceiver() (*Receiver, error) {
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

	return &Receiver{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}
func (r *Receiver) StartConsuming() error {
    msgs, err := r.channel.Consume(
        "Messages", 
        "",         
        false,      
        false,      
        false,      
        false,      
        nil,
    )
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
        for msg := range msgs {
				var message models.RequestBody
				if err := json.Unmarshal(msg.Body, &message); err != nil {
					// return err
					fmt.Printf("Error processing message: %v\n", err)
                msg.Nack(false, true)
                continue
				}
			
				log.Println("Saving to database...")
				query := `INSERT INTO messages (senderID, receiverID, content) VALUES ($1, $2, $3)`
				_, err := db.Exec(query, message.SenderID, message.ReceiverID, message.Content)
                if err != nil {
					fmt.Printf("Error saving to database: %v\n", err)
					msg.Nack(false, true)
					continue
				}
				msg.Ack(false)
        }
    }()

	return nil
}

func (r *Receiver) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}