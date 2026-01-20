package repository

import (
	"context"
	"fmt"
	"order-service/model"
	"time"

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
	GetSellerStatistics(sellerID string, from, to time.Time, groupBy string) (int, float64, []map[string]interface{}, error)
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

func (r *orderRepository) GetSellerStatistics(sellerID string, from, to time.Time, groupBy string) (int, float64, []map[string]interface{}, error) {
	// Build filter for seller and time range
	filter := bson.M{
		"seller._id": sellerID,
		"created_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	var pipeline []bson.M

	if groupBy == "" {
		// Overall statistics (no grouping)
		pipeline = []bson.M{
			{
				"$match": filter,
			},
			{
				"$group": bson.M{
					"_id": nil,
					"count": bson.M{"$sum": 1},
					"revenue": bson.M{"$sum": "$total"},
				},
			},
		}
	} else {
		// Breakdown by period (day or month)
		var dateFormat string
		if groupBy == "day" {
			dateFormat = "%Y-%m-%d"
		} else if groupBy == "month" {
			dateFormat = "%Y-%m"
		} else {
			return 0, 0, nil, fmt.Errorf("invalid groupBy parameter: %s", groupBy)
		}

		pipeline = []bson.M{
			{
				"$match": filter,
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"$dateToString": bson.M{
							"format": dateFormat,
							"date":   "$created_at",
						},
					},
					"count": bson.M{"$sum": 1},
					"revenue": bson.M{"$sum": "$total"},
				},
			},
			{
				"$sort": bson.M{"_id": 1}, // Sort by period ascending
			},
		}
	}

	cursor, err := r.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, 0, nil, err
	}
	defer cursor.Close(context.Background())

	// Decode results
	var results []map[string]interface{}
	if err = cursor.All(context.Background(), &results); err != nil {
		return 0, 0, nil, err
	}

	// If no orders found, return 0 for both
	if len(results) == 0 {
		return 0, 0, nil, nil
	}

	if groupBy == "" {
		// Overall statistics - return first (and only) result
		count := int(results[0]["count"].(int32))
		revenue := results[0]["revenue"].(float64)
		return count, revenue, nil, nil
	}

	// Calculate totals from breakdown
	totalCount := 0
	totalRevenue := 0.0
	for _, result := range results {
		totalCount += int(result["count"].(int32))
		totalRevenue += result["revenue"].(float64)
	}

	return totalCount, totalRevenue, results, nil
}

