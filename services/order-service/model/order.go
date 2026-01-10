package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              string        `bson:"_id" json:"id"`
	Status          string        `bson:"status" json:"status"`
	User            User          `bson:"user" json:"user"`
	PaymentMethod   string        `bson:"payment_method" json:"payment_method"`      //COD, STRIPE
	PaymentStatus   string        `bson:"payment_status" json:"payment_status"`      //PENDING, PAID, FAILED
	Seller          User          `bson:"seller" json:"seller"`
	Items           []OrderItem   `bson:"items" json:"items"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
	Voucher         *OrderVoucher `bson:"voucher" json:"voucher"`
	Total           float64       `bson:"total" json:"total"`
	Phone           string        `bson:"phone" json:"phone"`
	ShippingAddress OrderAddress  `bson:"shipping_address" json:"shipping_address"`
	DeliveryCode    string        `bson:"delivery_code" json:"delivery_code"`
}

type OrderItem struct {
	ProductID   string `bson:"product_id" json:"product_id"`
	VariantID   string `bson:"variant_id" json:"variant_id"`
	ProductName string `bson:"product_name" json:"product_name"`
	VariantName string `bson:"variant_name" json:"variant_name"` // e.g., "Size: M, Color: Red"
	SKU         string `bson:"sku" json:"sku"`
	Price       int    `bson:"price" json:"price"` // Price at the time of purchase
	Image       string `bson:"image" json:"image"`
	Quantity    int    `bson:"quantity" json:"quantity"`
}

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
}

type OrderAddress struct {
	FullName    string  `bson:"full_name" json:"full_name"`
	Phone       string  `bson:"phone" json:"phone"`
	AddressLine string  `bson:"address_line" json:"address_line"`
	Ward        string  `bson:"ward" json:"ward"`
	District    string  `bson:"district" json:"district"`
	Province    string  `bson:"province" json:"province"`
	Country     string  `bson:"country" json:"country"`
	Latitude    float64 `bson:"latitude" json:"latitude"`
	Longitude   float64 `bson:"longitude" json:"longitude"`
}

type OrderVoucher struct {
	Code                   string   `bson:"code" json:"code"`
	DiscountType           string   `bson:"discount_type" json:"discount_type"`
	DiscountValue          int      `bson:"discount_value" json:"discount_value"`
	MaxDiscountValue       *int     `bson:"max_discount_value" json:"max_discount_value"`
	MinOrderValue          int      `bson:"min_order_value" json:"min_order_value"`
	ApplyScope             string   `bson:"apply_scope" json:"apply_scope"`
	ApplySellerCategoryIds []string `bson:"apply_seller_category_ids" json:"apply_seller_category_ids"`
}

func (o *Order) BeforeCreate() {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = time.Now()
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = time.Now()
	}
	if o.Status == "" {
		o.Status = "pending"
	}
}