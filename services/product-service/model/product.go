package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `bson:"_id" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Images      []ProductImages `bson:"images" json:"images"`
	Status      string    `bson:"status" json:"status"`
	SellerID    string    `bson:"seller_id" json:"seller_id"`
	Rating      float64   `bson:"rating" json:"rating"`
	RateCount   int       `bson:"rate_count" json:"rate_count"`
	SoldCount   int       `bson:"sold_count" json:"sold_count"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Price       Price     `bson:"price" json:"price"`
	Stock       int       `bson:"stock" json:"stock"`

	OptionGroups []OptionGroup      `bson:"option_groups" json:"option_groups"` 
	Variants     []Variant          `bson:"variants" json:"variants"`

	// Many-to-many relationships stored as arrays of IDs
	CategoryIDs       []string `bson:"category_ids" json:"category_ids,omitempty"`
	SellerCategoryIDs []string `bson:"seller_category_ids" json:"seller_category_ids,omitempty"`
}

type Price struct {
	Min int `bson:"min" json:"min"`
	Max int `bson:"max" json:"max"`
}

type OptionGroup struct {
	Key    string   `bson:"key"`   // size, color
	Values []string `bson:"values"`
}

type Variant struct {
	ID      string            `bson:"_id" json:"id"`
	SKU     string            `bson:"sku" json:"sku"`
	Options map[string]string `bson:"options" json:"options"`
	Price   int               `bson:"price" json:"price"`
	Stock   int               `bson:"stock" json:"stock"`
	Image   string            `bson:"image" json:"image"`
	SoldCount int               `bson:"sold_count" json:"sold_count"`
}

// Validate checks if the product has all required fields and valid data
func (p *Product) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if p.SellerID == "" {
		return fmt.Errorf("seller_id is required")
	}
	if len(p.Variants) == 0 {
		return fmt.Errorf("at least one variant is required")
	}

	// Validate each variant
	for i, variant := range p.Variants {
		if variant.SKU == "" {
			return fmt.Errorf("variant %d: SKU is required", i)
		}
		if variant.Price <= 0 {
			return fmt.Errorf("variant %d: price must be greater than 0", i)
		}
		if variant.Stock < 0 {
			return fmt.Errorf("variant %d: stock cannot be negative", i)
		}

		// Validate that variant options match option groups
		if len(p.OptionGroups) > 0 {
			for _, optionGroup := range p.OptionGroups {
				optionValue, exists := variant.Options[optionGroup.Key]
				if !exists {
					return fmt.Errorf("variant %d: missing option for '%s'", i, optionGroup.Key)
				}
				// Check if the option value is valid
				validValue := false
				for _, validVal := range optionGroup.Values {
					if validVal == optionValue {
						validValue = true
						break
					}
				}
				if !validValue {
					return fmt.Errorf("variant %d: invalid value '%s' for option '%s'", i, optionValue, optionGroup.Key)
				}
			}
		}
	}

	return nil
}

// BeforeCreate generates a new UUID for the ID field if not set and initializes defaults
func (p *Product) BeforeCreate() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}

	// Initialize nested ProductImages IDs
	for i := range p.Images {
		p.Images[i].BeforeCreate()
	}

	// Initialize empty slices if nil
	if p.CategoryIDs == nil {
		p.CategoryIDs = []string{}
	}
	if p.SellerCategoryIDs == nil {
		p.SellerCategoryIDs = []string{}
	}
	if p.Images == nil {
		p.Images = []ProductImages{}
	}
	if p.OptionGroups == nil {
		p.OptionGroups = []OptionGroup{}
	}
	if p.Variants == nil {
		p.Variants = []Variant{}
	}

	// Initialize variant IDs and sold counts
	for i := range p.Variants {
		if p.Variants[i].ID == "" {
			p.Variants[i].ID = uuid.New().String()
		}
		// Initialize variant sold count to 0
		p.Variants[i].SoldCount = 0
	}

	// Initialize product sold count to 0
	p.SoldCount = 0

	// Calculate stock, price (min/max) based on variants
	if len(p.Variants) > 0 {
		totalStock := 0
		minPrice := p.Variants[0].Price
		maxPrice := p.Variants[0].Price

		for _, variant := range p.Variants {
			totalStock += variant.Stock
			
			if variant.Price < minPrice {
				minPrice = variant.Price
			}
			if variant.Price > maxPrice {
				maxPrice = variant.Price
			}
		}

		// Set price range
		p.Price.Min = minPrice
		p.Price.Max = maxPrice

		// Set stock
		p.Stock = totalStock

		// Set status based on stock availability
		if p.Status == "" {
			if totalStock > 0 {
				p.Status = "available"
			} else {
				p.Status = "out_of_stock"
			}
		}
	} else {
		// No variants: set default status
		if p.Status == "" {
			p.Status = "draft"
		}
	}
}

// BeforeUpdate updates the UpdatedAt timestamp
func (p *Product) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}
