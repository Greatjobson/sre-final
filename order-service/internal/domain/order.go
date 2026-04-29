package domain

type Order struct {
	ID     string
	UserID string
	Status string
	Total  float64
}
