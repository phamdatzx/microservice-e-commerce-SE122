package repository

import (
	"context"
	"order-service/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
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

func (r *orderRepository) CreateOrder(order *model.Order) error {
	order.BeforeCreate()
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}
