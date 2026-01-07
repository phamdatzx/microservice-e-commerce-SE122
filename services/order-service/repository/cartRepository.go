package repository

import (
	"context"
	"order-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	FindCartItemByUserAndVariant(userID, variantID string) (*model.CartItem, error)
	CreateCartItem(item *model.CartItem) error
	UpdateCartItemQuantity(id string, quantity int) error
	FindCartItemByID(id string) (*model.CartItem, error)
	DeleteCartItem(id string) error
	FindCartItemsByUser(userID string) ([]model.CartItem, error)
	GetCartItemCount(userID string) (int64, error)
}

type cartRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewCartRepository(db *mongo.Database) CartRepository {
	return &cartRepository{
		db:         db,
		collection: db.Collection("cart_items"),
	}
}

func (r *cartRepository) FindCartItemByUserAndVariant(userID, variantID string) (*model.CartItem, error) {
	var cartItem model.CartItem
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":    userID,
		"variant.id": variantID,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&cartItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No existing cart item found
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepository) CreateCartItem(item *model.CartItem) error {
	item.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *cartRepository) UpdateCartItemQuantity(id string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"quantity":   quantity,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *cartRepository) FindCartItemByID(id string) (*model.CartItem, error) {
	var cartItem model.CartItem
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&cartItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Cart item not found
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepository) DeleteCartItem(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
func (r *cartRepository) FindCartItemsByUser(userID string) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &cartItems); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (r *cartRepository) GetCartItemCount(userID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$quantity"}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return 0, err
	}

	if len(result) > 0 {
		// Extract total from result
		total, ok := result[0]["total"].(int32)
		if ok {
			return int64(total), nil
		}
		// Handle if it's int64 or other number type just in case, though mongo driver usually decodes to int32/int64
		total64, ok := result[0]["total"].(int64)
		if ok {
			return total64, nil
		}
	}

	return 0, nil
}
