package repository

import (
	"context"
	"order-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	FindOrderByID(id string) (*model.Order, error)
	UpdateOrder(order *model.Order) error
	FindOrdersByUser(userID string, status string, page, limit int, sortBy, sortOrder string) ([]*model.Order, int64, error)
	FindOrdersBySeller(sellerID string, status string, paymentMethod string, paymentStatus string, search string, page, limit int, sortBy, sortOrder string) ([]*model.Order, int64, error)
	VerifyVariantPurchase(userID, productID, variantID string) (bool, error)
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

func (r *orderRepository) FindOrdersByUser(userID string, status string, page, limit int, sortBy, sortOrder string) ([]*model.Order, int64, error) {
	// Build filter
	filter := bson.M{"user._id": userID}
	if status != "" {
		filter["status"] = status
	}

	// Count total documents
	totalCount, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortOptions := bson.D{}
	sortDirection := -1 // default descending
	if sortOrder == "asc" {
		sortDirection = 1
	}

	switch sortBy {
	case "total":
		sortOptions = bson.D{{Key: "total", Value: sortDirection}}
	case "created_at":
		sortOptions = bson.D{{Key: "created_at", Value: sortDirection}}
	default:
		sortOptions = bson.D{{Key: "created_at", Value: -1}} // default: newest first
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Build find options
	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Decode results
	var orders []*model.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (r *orderRepository) FindOrdersBySeller(sellerID string, status string, paymentMethod string, paymentStatus string, search string, page, limit int, sortBy, sortOrder string) ([]*model.Order, int64, error) {
	// Build filter
	filter := bson.M{"seller._id": sellerID}
	
	if status != "" {
		filter["status"] = status
	}
	
	if paymentMethod != "" {
		filter["payment_method"] = paymentMethod
	}
	
	if paymentStatus != "" {
		filter["payment_status"] = paymentStatus
	}
	
	// Search by order ID or phone
	if search != "" {
		filter["$or"] = []bson.M{
			{"_id": bson.M{"$regex": search, "$options": "i"}},
			{"phone": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Count total documents
	totalCount, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortOptions := bson.D{}
	sortDirection := -1 // default descending
	if sortOrder == "asc" {
		sortDirection = 1
	}

	switch sortBy {
	case "total":
		sortOptions = bson.D{{Key: "total", Value: sortDirection}}
	case "created_at":
		sortOptions = bson.D{{Key: "created_at", Value: sortDirection}}
	default:
		sortOptions = bson.D{{Key: "created_at", Value: -1}} // default: newest first
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Build find options
	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	// Decode results
	var orders []*model.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (r *orderRepository) VerifyVariantPurchase(userID, productID, variantID string) (bool, error) {
	filter := bson.M{
		"user._id": userID,
		"items": bson.M{
			"$elemMatch": bson.M{
				"product_id": productID,
				"variant_id": variantID,
			},
		},
	}

	count, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

