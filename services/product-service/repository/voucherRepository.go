package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VoucherRepository interface {
	Create(voucher *model.Voucher) error
	FindByID(id string) (*model.Voucher, error)
	FindBySellerID(sellerID string) ([]model.Voucher, error)
	Update(voucher *model.Voucher) error
	Delete(id string) error
}

type voucherRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewVoucherRepository(db *mongo.Database) VoucherRepository {
	return &voucherRepository{
		db:         db,
		collection: db.Collection("vouchers"),
	}
}

func (r *voucherRepository) Create(voucher *model.Voucher) error {
	voucher.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, voucher)
	return err
}

func (r *voucherRepository) FindByID(id string) (*model.Voucher, error) {
	var voucher model.Voucher
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&voucher)
	return &voucher, err
}

func (r *voucherRepository) FindBySellerID(sellerID string) ([]model.Voucher, error) {
	var vouchers []model.Voucher
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{"seller_id": sellerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &vouchers); err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (r *voucherRepository) Update(voucher *model.Voucher) error {
	voucher.BeforeUpdate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": voucher.ID},
		bson.M{"$set": voucher},
	)
	return err
}

func (r *voucherRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
