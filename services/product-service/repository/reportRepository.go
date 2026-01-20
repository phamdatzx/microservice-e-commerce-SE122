package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportRepository interface {
	Create(report *model.Report) error
	FindByID(id string) (*model.Report, error)
	FindAll(skip, limit int, status *model.ReportStatus, sortOrder int) ([]model.Report, int64, error)
	FindByProductID(productID string, skip, limit int) ([]model.Report, int64, error)
	FindByUserID(userID string, skip, limit int) ([]model.Report, int64, error)
	Update(report *model.Report) error
	Delete(id string) error
	FindByProductIDAndUserID(productID, userID string) (*model.Report, error)
}

type reportRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewReportRepository(db *mongo.Database) ReportRepository {
	return &reportRepository{
		db:         db,
		collection: db.Collection("reports"),
	}
}

func (r *reportRepository) Create(report *model.Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, report)
	return err
}

func (r *reportRepository) FindByID(id string) (*model.Report, error) {
	var report model.Report
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&report)
	return &report, err
}

func (r *reportRepository) FindAll(skip, limit int, status *model.ReportStatus, sortOrder int) ([]model.Report, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	// Add status filter if provided
	if status != nil {
		filter["status"] = *status
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination, sorted by created_at
	// sortOrder: 1 for ascending, -1 for descending
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: sortOrder}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []model.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) FindByProductID(productID string, skip, limit int) ([]model.Report, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"product_id": productID}

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

	var reports []model.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) FindByUserID(userID string, skip, limit int) ([]model.Report, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user._id": userID}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []model.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) Update(report *model.Report) error {
	report.UpdatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": report.ID},
		bson.M{"$set": report},
	)
	return err
}

func (r *reportRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *reportRepository) FindByProductIDAndUserID(productID, userID string) (*model.Report, error) {
	var report model.Report
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"product_id": productID, "user._id": userID}).Decode(&report)
	return &report, err
}
