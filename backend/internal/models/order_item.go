package models

type OrderItem struct {
	ID            string  `json:"id" bson:"_id,omitempty"`
	OrderID       string  `json:"order_id" bson:"orderId"`
	ProductID     string  `json:"product_id" bson:"productId"`
	ProductName   string  `json:"product_name" bson:"productName"`
	SelectedSize  string  `json:"selected_size" bson:"selectedSize"`
	SelectedColor string  `json:"selected_color" bson:"selectedColor"`
	Quantity      int     `json:"quantity" bson:"quantity"`
	UnitPrice     float64 `json:"unit_price" bson:"unitPrice"`
	LineTotal     float64 `json:"line_total" bson:"lineTotal"`
}
