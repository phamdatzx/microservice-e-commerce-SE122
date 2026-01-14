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
	FindBySeller(sellerID string, filter bson.M, skip, limit int, sortField string, sortDirection int) ([]model.Product, int64, error)
	FindVariantsByIds(variantIDs []string) (map[string]*model.Product, error)
	UpdateVariantStock(productID string, variantID string, stockDelta int) error
	SearchProducts(filter bson.M, skip, limit int, sortByTextScore bool) ([]model.Product, int64, error)
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

func (r *productRepository) FindBySeller(sellerID string, filter bson.M, skip, limit int, sortField string, sortDirection int) ([]model.Product, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add seller_id to filter
	filter["seller_id"] = sellerID

	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build find options with pagination and sorting
	opts := make([]interface{}, 0)
	
	// Add match stage
	opts = append(opts, bson.M{"$match": filter})
	
	// Add sort stage if sortField is provided
	if sortField != "" {
		opts = append(opts, bson.M{"$sort": bson.M{sortField: sortDirection}})
	}
	
	// Add skip and limit stages
	opts = append(opts, bson.M{"$skip": skip})
	opts = append(opts, bson.M{"$limit": limit})

	// Execute aggregation
	cursor, err := r.collection.Aggregate(ctx, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindVariantsByIds(variantIDs []string) (map[string]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query products that contain any of the variant IDs
	filter := bson.M{
		"variants._id": bson.M{"$in": variantIDs},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	// Create a map of variant_id -> product for quick lookup
	variantToProduct := make(map[string]*model.Product)
	for i := range products {
		for j := range products[i].Variants {
			variantID := products[i].Variants[j].ID
			// Check if this variant is in the requested list
			for _, requestedID := range variantIDs {
				if variantID == requestedID {
					variantToProduct[variantID] = &products[i]
					break
				}
			}
		}
	}

	return variantToProduct, nil
}

func (r *productRepository) UpdateVariantStock(productID string, variantID string, stockDelta int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the product first to find the variant
	var product model.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return err
	}

	// Find the variant and update its stock
	variantFound := false
	for i := range product.Variants {
		if product.Variants[i].ID == variantID {
			product.Variants[i].Stock += stockDelta
			variantFound = true
			break
		}
	}

	if !variantFound {
		return mongo.ErrNoDocuments
	}

	// Update the entire variants array, product stock, and updated_at
	product.UpdatedAt = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": productID},
		bson.M{
			"$set": bson.M{
				"variants":   product.Variants,
				"updated_at": product.UpdatedAt,
			},
			"$inc": bson.M{
				"stock": stockDelta,
			},
		},
	)
	return err
}

func (r *productRepository) SearchProducts(filter bson.M, skip, limit int, sortByTextScore bool) ([]model.Product, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build aggregation pipeline
	pipeline := make([]bson.M, 0)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": filter})

	// Add text score projection and sort if text search is being used
	if sortByTextScore {
		// Add text score to sort by relevance
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"textScore": bson.M{"$meta": "textScore"},
			},
		})
		pipeline = append(pipeline, bson.M{
			"$sort": bson.M{"textScore": -1},
		})
	}

	// Add pagination
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	// Execute aggregation
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
