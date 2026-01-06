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
