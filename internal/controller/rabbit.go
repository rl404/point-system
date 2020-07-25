package controller

import (
	"encoding/json"
	"time"

	"github.com/rl404/point-system/internal/config"
	"github.com/streadway/amqp"
)

const (
	// ActionAdd is action to add point.
	ActionAdd = "add"
	// ActionSubtract is action to subtract point.
	ActionSubtract = "subtract"
)

// RabbitQueue is model data for rabbitmq queue.
type RabbitQueue struct {
	UserID      int       `json:"user_id"`
	Action      string    `json:"action"`
	RequestedAt time.Time `json:"requested_at"`
}

// sendQueue to send data to rabbitmq queue.
func (ph *PointHandler) sendQueue(data RabbitQueue) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ph.Rabbit.Publish(
		"",                  // exchange
		config.RabbitMQName, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType:  "text/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}
