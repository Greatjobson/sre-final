package domain

type User struct {
	ID           string
	FullName     string
	Email        string
	PasswordHash string
	Role         string
}
