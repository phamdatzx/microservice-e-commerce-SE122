package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	// TODO: Add order repository methods
}

type orderRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) OrderRepository {
	return &orderRepository{
		db:         db,
		collection: db.Collection("orders"),
	}
}

// TODO: Implement order repository methods
