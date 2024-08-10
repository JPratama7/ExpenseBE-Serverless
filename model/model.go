package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at"`
}

type Transaction struct {
	Id              primitive.ObjectID `bson:"_id"`
	UserId          primitive.ObjectID `bson:"user_id"`
	Description     string             `bson:"description"`
	Amount          float64            `bson:"amount"`
	TransactionDate primitive.DateTime `bson:"transaction_date"`
	Category        primitive.ObjectID `bson:"category"`
	IsIncome        bool               `bson:"is_income"`
	CreatedAt       primitive.DateTime `bson:"created_at"`
	UpdatedAt       primitive.DateTime `bson:"updated_at"`
}

type Categories struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Name      primitive.ObjectID `bson:"name"`
	Type      string             `bson:"type"`
	CreatedAt primitive.ObjectID `bson:"created_at"`
	UpdatedAt primitive.ObjectID `bson:"updated_at"`
}

type Budget struct {
	Id          primitive.ObjectID `bson:"_id"`
	UserId      primitive.ObjectID `bson:"user_id"`
	CategoryId  primitive.ObjectID `bson:"category_id"`
	LimitAmount float64            `bson:"limit_amount"`
	StartDate   primitive.DateTime `bson:"start_date"`
	EndDate     primitive.DateTime `bson:"end_date"`
	CreatedAt   primitive.DateTime `bson:"created_at"`
	UpdatedAt   primitive.DateTime `bson:"updated_at"`
}
