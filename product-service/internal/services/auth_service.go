package services

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/product-service/internal/models"
	"github.com/Tedra-ez/AdvancedProgramming_Final/product-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

type AuthService struct {
	users  *repository.UserRepository
	secret []byte
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	return &AuthService{
		users:  users,
		secret: []byte(secret),
	}
}

func (s *AuthService) Register(ctx context.Context, fullName, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		return ErrEmailExists
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	role := "customer"
	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" && strings.EqualFold(email, adminEmail) {
		role = "admin"
	}
	user := &models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}

	return s.users.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"role":  user.Role,
		"email": user.Email,
		"name":  user.FullName,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *AuthService) ParseToken(ctx context.Context, tokenStr string) (map[string]string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return map[string]string{
		"id":    getString(claims, "sub"),
		"role":  getString(claims, "role"),
		"email": getString(claims, "email"),
		"name":  getString(claims, "name"),
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, tokenStr string) (string, error) {
	data, err := s.ParseToken(ctx, tokenStr)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	user, err := s.users.FindByID(ctx, data["id"])
	if err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.Hex(),
		"role":  user.Role,
		"email": user.Email,
		"name":  user.FullName,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *AuthService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.users.FindAll(ctx)
}

func (s *AuthService) GetUserCount(ctx context.Context) (int64, error) {
	return s.users.Count(ctx)
}

func getString(m jwt.MapClaims, key string) string {
	if v, ok := m[key]; ok && v != nil {
		switch s := v.(type) {
		case string:
			return s
		default:
			return ""
		}
	}
	return ""
}
