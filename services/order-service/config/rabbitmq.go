package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

// InitRabbitMQ initializes a shared RabbitMQ connection and channel.
// If configuration is missing or connection fails, the service will continue running
// but message publishing will be disabled (errors will be logged).
func InitRabbitMQ() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		fmt.Println("⚠️  Missing RabbitMQ configuration (RABBITMQ_URL), interaction events will be disabled")
		return
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Printf("❌ Failed to connect to RabbitMQ: %v\n", err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("❌ Failed to open RabbitMQ channel: %v\n", err)
		_ = conn.Close()
		return
	}

	rabbitConn = conn
	rabbitChannel = ch

	fmt.Println("✅ RabbitMQ connected (order-service)")
}

// CloseRabbitMQ closes the shared RabbitMQ connection and channel.
func CloseRabbitMQ() {
	if rabbitChannel != nil {
		_ = rabbitChannel.Close()
	}
	if rabbitConn != nil {
		_ = rabbitConn.Close()
	}
}

// UserInteractionMessage is the payload published to the user interaction queue.
// It is consumed by the ai-service user_vector_worker.
type UserInteractionMessage struct {
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Action    string  `json:"action"`
	Score     float64 `json:"score"`
}

// PublishUserInteraction publishes a user–product interaction event to RabbitMQ.
// action should be one of: "view", "add_to_cart", "purchase".
// score is the weight for the action (view=1, add_to_cart=10, purchase=10).
// This function is best-effort: failures are logged but do not affect the main flow.
func PublishUserInteraction(userID, productID, action string, score float64) {
	if rabbitChannel == nil {
		fmt.Println("⚠️  RabbitMQ channel is nil, skipping user interaction publish")
		return
	}

	if userID == "" || productID == "" {
		fmt.Printf("⚠️  Skipping user interaction publish: userID=%q productID=%q\n", userID, productID)
		return
	}

	msg := UserInteractionMessage{
		UserID:    userID,
		ProductID: productID,
		Action:    action,
		Score:     score,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ Failed to marshal user interaction message: %v\n", err)
		return
	}

	queueName := os.Getenv("RABBITMQ_USER_INTERACTION_QUEUE")
	if queueName == "" {
		queueName = "user.interaction"
	}

	// Ensure queue exists (idempotent)
	_, err = rabbitChannel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		fmt.Printf("❌ Failed to declare RabbitMQ queue '%s': %v\n", queueName, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = rabbitChannel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key (queue)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		fmt.Printf("❌ Failed to publish user interaction message: %v\n", err)
		return
	}

	fmt.Printf("✅ Published user interaction [%s] userID=%s productID=%s score=%.0f to '%s'\n",
		action, userID, productID, score, queueName)
}
