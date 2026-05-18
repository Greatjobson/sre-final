package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName         string             `bson:"fullName" json:"fullName"`
	Email            string             `bson:"email" json:"email"`
	PasswordHash     string             `bson:"passwordHash" json:"-"`
	Role             string             `bson:"role" json:"role"`
	IsActive         bool               `bson:"isActive" json:"isActive"`
	TwoFactorEnabled bool               `bson:"twoFactorEnabled" json:"twoFactorEnabled"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
