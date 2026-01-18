package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchHistoryRepository interface {
	Create(searchHistory *model.SearchHistory) error
	FindByUserID(userID string, limit int) ([]model.SearchHistory, error)
}

type searchHistoryRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSearchHistoryRepository(db *mongo.Database) SearchHistoryRepository {
	return &searchHistoryRepository{
		db:         db,
		collection: db.Collection("search_histories"),
	}
}

func (r *searchHistoryRepository) Create(searchHistory *model.SearchHistory) error {
	searchHistory.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, searchHistory)
	return err
}

func (r *searchHistoryRepository) FindByUserID(userID string, limit int) ([]model.SearchHistory, error) {
	var searchHistories []model.SearchHistory
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Find options: sort by searched_at descending and limit results
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "searched_at", Value: -1}})
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &searchHistories); err != nil {
		return nil, err
	}
	return searchHistories, nil
}
