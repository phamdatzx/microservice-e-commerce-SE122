package model

import "time"

type Order struct {
	ID        string    `bson:"_id" json:"id"`
	Status    string    `bson:"status" json:"status"`
	User      User      `bson:"user" json:"user"`
	Seller    User      `bson:"seller" json:"seller"`
	Items     []OrderItem `bson:"items" json:"items"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	VoucherID string    `bson:"voucher_id" json:"voucher_id"`
	Total     float64   `bson:"total" json:"total"`
	Phone 	string 	
}

type OrderItem struct {
	ProductID string    `bson:"product_id" json:"product_id"`
	VariantID string    `bson:"variant_id" json:"variant_id"`
	Quantity  int       `bson:"quantity" json:"quantity"`
}

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
}