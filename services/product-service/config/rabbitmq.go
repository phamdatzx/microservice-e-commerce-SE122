package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"product-service/model"

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
		fmt.Println("⚠️  Missing RabbitMQ configuration (RABBITMQ_URL)")
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

	fmt.Println("✅ RabbitMQ connected")
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

type ProductCreatedMessage struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	SellerID     string `json:"seller_id"`
}

// PublishProductCreated publishes a product.created event to RabbitMQ.
// This function is best-effort: failures are logged but do not affect the main flow.
func PublishProductCreated(product *model.Product, categoryName string) {
	if rabbitChannel == nil {
		// RabbitMQ not initialized; nothing to do
		fmt.Println("⚠️  RabbitMQ channel is nil, skipping product.created publish")
		return
	}

	if len(product.CategoryIDs) == 0 {
		// No category to publish; skip
		fmt.Println("⚠️  Product has no category_ids, skipping product.created publish")
		return
	}

	categoryID := product.CategoryIDs[0]

	msg := ProductCreatedMessage{
		ID:           product.ID,
		Name:         product.Name,
		CategoryID:   categoryID,
		CategoryName: categoryName,
		SellerID:     product.SellerID,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ Failed to marshal product.created message: %v\n", err)
		return
	}

	queueName := os.Getenv("RABBITMQ_PRODUCT_CREATED_QUEUE")
	if queueName == "" {
		queueName = "product.created"
	}

	// Ensure queue exists
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
		fmt.Printf("❌ Failed to publish product.created message: %v\n", err)
		return
	}

	fmt.Printf("✅ Published product.created message to queue '%s'\n", queueName)
}

