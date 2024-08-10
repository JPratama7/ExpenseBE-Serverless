package repository

import "crud/model"

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
	Create(username, email, passwordHash string) (model.User, error)
	GetAll() ([]model.User, error)
	GetById(id string) (model.User, error)
	GetByUsername(username string) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetByEmailOrUsername(email, userName string) (model.User, error)
}
