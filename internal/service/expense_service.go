package service

import (
	"context"
	"time"

	"expense-tracker-pwa/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExpenseService struct {
	expenses *mongo.Collection
}

type CreateExpenseInput struct {
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Note     string    `json:"note"`
	Date     time.Time `json:"date"`
}

func NewExpenseService(db *mongo.Database) *ExpenseService {
	return &ExpenseService{expenses: db.Collection("expenses")}
}

func (s *ExpenseService) List(ctx context.Context, userID primitive.ObjectID) ([]model.Expense, error) {
	cur, err := s.expenses.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var items []model.Expense
	if err := cur.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ExpenseService) Create(ctx context.Context, userID primitive.ObjectID, in CreateExpenseInput) (model.Expense, error) {
	e := model.Expense{UserID: userID, Amount: in.Amount, Category: in.Category, Note: in.Note, Date: in.Date, CreatedAt: time.Now()}
	res, err := s.expenses.InsertOne(ctx, e)
	if err != nil {
		return model.Expense{}, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	e.ID = id
	return e, nil
}
