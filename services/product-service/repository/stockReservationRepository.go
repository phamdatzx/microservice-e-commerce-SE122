package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StockReservationRepository interface {
	Create(reservation *model.StockReservation) error
	CreateMany(reservations []model.StockReservation) error
	FindByOrderID(orderID string) ([]model.StockReservation, error)
	UpdateStatusByOrderID(orderID string, status string) error
}

type stockReservationRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewStockReservationRepository(db *mongo.Database) StockReservationRepository {
	return &stockReservationRepository{
		db:         db,
		collection: db.Collection("stock_reservations"),
	}
}

func (r *stockReservationRepository) Create(reservation *model.StockReservation) error {
	reservation.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, reservation)
	return err
}

func (r *stockReservationRepository) CreateMany(reservations []model.StockReservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize each reservation
	for i := range reservations {
		reservations[i].BeforeCreate()
	}

	// Convert to []interface{} for InsertMany
	docs := make([]interface{}, len(reservations))
	for i, reservation := range reservations {
		docs[i] = reservation
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *stockReservationRepository) FindByOrderID(orderID string) ([]model.StockReservation, error) {
	var reservations []model.StockReservation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *stockReservationRepository) UpdateStatusByOrderID(orderID string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"order_id": orderID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}
