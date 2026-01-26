package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VoucherUsageRepository interface {
	Create(usage *model.VoucherUsage) error
	CountByUserAndVoucher(userID, voucherID string) (int64, error)
	FindByUserAndVoucher(userID, voucherID string) ([]model.VoucherUsage, error)
}

type voucherUsageRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewVoucherUsageRepository(db *mongo.Database) VoucherUsageRepository {
	return &voucherUsageRepository{
		db:         db,
		collection: db.Collection("voucher_usages"),
	}
}

func (r *voucherUsageRepository) Create(usage *model.VoucherUsage) error {
	usage.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, usage)
	return err
}

func (r *voucherUsageRepository) CountByUserAndVoucher(userID, voucherID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"user_id":    userID,
		"voucher_id": voucherID,
	})
	return count, err
}

func (r *voucherUsageRepository) FindByUserAndVoucher(userID, voucherID string) ([]model.VoucherUsage, error) {
	var usages []model.VoucherUsage
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":    userID,
		"voucher_id": voucherID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &usages); err != nil {
		return nil, err
	}
	return usages, nil
}
