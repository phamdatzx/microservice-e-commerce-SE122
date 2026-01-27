package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SavedVoucherRepository interface {
	Create(savedVoucher *model.SavedVoucher) error
	FindByUserID(userID string) ([]model.SavedVoucher, error)
	FindByUserAndVoucher(userID, voucherID string) (*model.SavedVoucher, error)
	Delete(userID, voucherID string) error
	DeleteByVoucherID(voucherID string) error
	Update(savedVoucher *model.SavedVoucher) error
}

type savedVoucherRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewSavedVoucherRepository(db *mongo.Database) SavedVoucherRepository {
	return &savedVoucherRepository{
		db:         db,
		collection: db.Collection("saved_vouchers"),
	}
}

func (r *savedVoucherRepository) Create(savedVoucher *model.SavedVoucher) error {
	savedVoucher.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, savedVoucher)
	return err
}

func (r *savedVoucherRepository) FindByUserID(userID string) ([]model.SavedVoucher, error) {
	var savedVouchers []model.SavedVoucher
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Only get non-deleted saved vouchers
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID, "is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &savedVouchers); err != nil {
		return nil, err
	}
	return savedVouchers, nil
}

func (r *savedVoucherRepository) FindByUserAndVoucher(userID, voucherID string) (*model.SavedVoucher, error) {
	var savedVoucher model.SavedVoucher
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Only get non-deleted saved vouchers
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID, "voucher_id": voucherID, "is_deleted": false}).Decode(&savedVoucher)
	return &savedVoucher, err
}

func (r *savedVoucherRepository) Delete(userID, voucherID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Soft delete: set is_deleted to true instead of removing the record
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID, "voucher_id": voucherID},
		bson.M{"$set": bson.M{"is_deleted": true}},
	)
	return err
}

func (r *savedVoucherRepository) DeleteByVoucherID(voucherID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Soft delete: set is_deleted to true for all saved vouchers with this voucher_id
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"voucher_id": voucherID},
		bson.M{"$set": bson.M{"is_deleted": true}},
	)
	return err
}

func (r *savedVoucherRepository) Update(savedVoucher *model.SavedVoucher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": savedVoucher.UserID, "voucher_id": savedVoucher.VoucherID},
		bson.M{"$set": bson.M{"used_count": savedVoucher.UsedCount}},
	)
	return err
}
