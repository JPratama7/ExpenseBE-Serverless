package repository

import (
	"context"
	"crud/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TransactionRepository interface {
	Create(userId, description, amount, transactionDate string, categoryId string, isIncome bool) error
	GetAll(userId string) ([]model.Transaction, error)
	GetById(id string) (model.Transaction, error)
}

type CategoryRepository interface {
	Create(userId, name, type_ string) error
	GetAll(userId string) ([]model.Categories, error)
	GetById(id string) (model.Categories, error)
}

type BudgetRepository interface {
	Create(userId, categoryId, limitAmount, startDate, endDate string) error
	GetAll(userId string) ([]model.Budget, error)
	GetById(id string) (model.Budget, error)
}

type UserRepository interface {
	Create(ctx context.Context, data model.User) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetById(ctx context.Context, id string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByEmailOrUsername(ctx context.Context, email, userName string) (model.User, error)
}

type userRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database, coll string) UserRepository {
	return &userRepository{
		coll: db.Collection(coll),
	}
}

func (r *userRepository) Create(ctx context.Context, data model.User) (res model.User, err error) {
	res = data
	res.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	res.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	cur, err := r.coll.InsertOne(ctx, res)
	if err != nil {
		return
	}
	_id, ok := cur.InsertedID.(primitive.ObjectID)
	if !ok {
		return
	}

	res.Id = _id
	return
}

func (r *userRepository) GetAll(ctx context.Context) (res []model.User, err error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}

func (r *userRepository) GetById(ctx context.Context, id string) (res model.User, err error) {
	cur := r.coll.FindOne(ctx, bson.M{"_id": id})
	if err = cur.Err(); err != nil {
		return
	}

	err = cur.Decode(&res)
	return
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (res model.User, err error) {
	cur := r.coll.FindOne(ctx, bson.M{"username": username})
	if err = cur.Err(); err != nil {
		return
	}

	err = cur.Decode(&res)
	return
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (res model.User, err error) {
	cur := r.coll.FindOne(ctx, bson.M{"email": email})
	if err = cur.Err(); err != nil {
		return
	}

	err = cur.Decode(&res)
	return
}

func (r *userRepository) GetByEmailOrUsername(ctx context.Context, email, userName string) (res model.User, err error) {
	cur := r.coll.FindOne(ctx, bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "email", Value: email}}, bson.D{{Key: "username", Value: userName}}}}})
	if err = cur.Err(); err != nil {
		return
	}
	err = cur.Decode(&res)
	return
}
