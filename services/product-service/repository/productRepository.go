package repository

import (
	"context"
	"product-service/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	FindAll() ([]model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
	AddImagesToProduct(productID string, images []model.ProductImages) error
	UpdateVariantImages(productID string, variantUpdates map[string]string) error
}

type productRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ProductRepository {
	return &productRepository{
		db:         db,
		collection: db.Collection("products"),
	}
}

func (r *productRepository) Create(product *model.Product) error {
	product.BeforeCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return &product, err
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(product *model.Product) error {
	product.BeforeUpdate()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": product.ID},
		bson.M{"$set": product},
	)
	return err
}

func (r *productRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *productRepository) AddImagesToProduct(productID string, images []model.ProductImages) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productID},
		bson.M{
			"$push": bson.M{
				"images": bson.M{"$each": images},
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

func (r *productRepository) UpdateVariantImages(productID string, variantUpdates map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the product first
	var product model.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return err
	}

	// Update variant images in memory
	for i := range product.Variants {
		if imageURL, exists := variantUpdates[product.Variants[i].ID]; exists {
			product.Variants[i].Image = imageURL
		}
	}

	// Update the entire variants array and updated_at
	product.UpdatedAt = time.Now()
	
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productID},
		bson.M{
			"$set": bson.M{
				"variants":   product.Variants,
				"updated_at": product.UpdatedAt,
			},
		},
	)
	return err
}
