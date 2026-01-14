package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RatingRepository interface {
	Create(rating *model.Rating) error
	FindByID(id string) (*model.Rating, error)
	FindAll() ([]model.Rating, error)
	FindByProductID(productID string, skip, limit int, star *int, hasImage *bool) ([]model.Rating, int64, error)
	FindByUserID(userID string) ([]model.Rating, error)
	Update(rating *model.Rating) error
	Delete(id string) error
	AddRatingResponse(ratingID string, response model.RatingResponse) error
	FindByProductIDAndUserID(productID string, userID string) (*model.Rating, error)
}

type ratingRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRatingRepository(db *mongo.Database) RatingRepository {
	return &ratingRepository{
		db:         db,
		collection: db.Collection("ratings"),
	}
}

func (r *ratingRepository) Create(rating *model.Rating) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, rating)
	return err
}

func (r *ratingRepository) FindByID(id string) (*model.Rating, error) {
	var rating model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rating)
	return &rating, err
}

func (r *ratingRepository) FindAll() ([]model.Rating, error) {
	var ratings []model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) FindByProductID(productID string, skip, limit int, star *int, hasImage *bool) ([]model.Rating, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"product_id": productID}

	// Add star filter if provided
	if star != nil {
		filter["star"] = *star
	}

	// Add hasImage filter if provided
	if hasImage != nil && *hasImage {
		// Filter ratings that have at least one image
		filter["images"] = bson.M{"$exists": true, "$ne": []interface{}{}}
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination, sorted by created_at descending
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var ratings []model.Rating
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, 0, err
	}

	return ratings, total, nil
}

func (r *ratingRepository) FindByUserID(userID string) ([]model.Rating, error) {
	var ratings []model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *ratingRepository) Update(rating *model.Rating) error {
	rating.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": rating.ID},
		bson.M{"$set": rating},
	)
	return err
}

func (r *ratingRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ratingRepository) AddRatingResponse(ratingID string, response model.RatingResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": ratingID},
		bson.M{
			"$push": bson.M{
				"rating_response": response,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *ratingRepository) FindByProductIDAndUserID(productID string, userID string) (*model.Rating, error) {
	var rating model.Rating
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"product_id": productID, "user._id": userID}).Decode(&rating)
	return &rating, err
}