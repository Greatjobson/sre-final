package models

import "time"

type Order struct {
	ID              string      `json:"id" bson:"_id,omitempty"`
	UserID          string      `json:"user_id" bson:"userId"`
	Status          string      `json:"status" bson:"status"`
	PaymentMethod   string      `json:"payment_method" bson:"paymentMethod"`
	DeliveryMethod  string      `json:"delivery_method" bson:"deliveryMethod"`
	DeliveryAddress string      `json:"delivery_address" bson:"deliveryAddress"`
	Comment         string      `json:"comment" bson:"comment"`
	Subtotal        float64     `json:"subtotal" bson:"subtotal"`
	DeliveryFee     float64     `json:"delivery_fee" bson:"deliveryFee"`
	Total           float64     `json:"total" bson:"total"`
	Items           []OrderItem `json:"items" bson:"-"`
	CreatedAt       time.Time   `json:"created_at" bson:"createdAt"`
	UpdatedAt       time.Time   `json:"updated_at" bson:"updatedAt"`
}

type CreateOrderRequest struct {
	UserID          string            `json:"user_id" binding:"required"`
	PaymentMethod   string            `json:"payment_method"`
	DeliveryMethod  string            `json:"delivery_method"`
	DeliveryAddress string            `json:"delivery_address"`
	Comment         string            `json:"comment"`
	Items           []CreateOrderItem `json:"items" binding:"required"`
}

type CreateOrderItem struct {
	ProductID     string  `json:"product_id" binding:"required"`
	ProductName   string  `json:"product_name"`
	SelectedSize  string  `json:"selected_size"`
	SelectedColor string  `json:"selected_color"`
	Quantity      int     `json:"quantity" binding:"required"`
	UnitPrice     float64 `json:"unit_price" binding:"required"`
}
