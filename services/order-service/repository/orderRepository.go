package repository

import (
	"context"
	"order-service/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	FindOrderByID(id string) (*model.Order, error)
	UpdateOrder(order *model.Order) error
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

func (r *orderRepository) FindOrderByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.collection.FindOne(context.Background(), map[string]string{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrder(order *model.Order) error {
	_, err := r.collection.ReplaceOne(context.Background(), map[string]string{"_id": order.ID}, order)
	return err
}
