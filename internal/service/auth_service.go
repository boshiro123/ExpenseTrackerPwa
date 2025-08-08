package service

import (
	"context"
	"errors"
	"time"

	"expense-tracker-pwa/internal/config"
	"expense-tracker-pwa/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     *mongo.Collection
	jwtSecret string
}

type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}

func NewAuthService(db *mongo.Database, cfg config.Config) *AuthService {
	return &AuthService{users: db.Collection("users"), jwtSecret: cfg.JWTSecret}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	var existing model.User
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&existing)
	if err == nil {
		return errors.New("user_exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := model.User{Email: email, PasswordHash: string(hash), CreatedAt: time.Now()}
	_, err = s.users.InsertOne(ctx, user)
	return err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user model.User
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid_credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", errors.New("invalid_credentials")
	}
	claims := jwt.MapClaims{
		"sub": user.ID.Hex(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseUserIDFromToken(tokenString, secret string) (primitive.ObjectID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
	if err != nil || !token.Valid {
		return primitive.NilObjectID, errors.New("invalid_token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid_token")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid_token")
	}
	id, err := primitive.ObjectIDFromHex(sub)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid_token")
	}
	return id, nil
}
