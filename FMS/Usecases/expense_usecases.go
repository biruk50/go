package Usecases

import (
	"FMS/Domain"
	"FMS/Repositories"
	"errors"
	"time"
)

type ExpenseUsecase interface {
	CreateExpense(input *Domain.Expense) (*Domain.Expense, error)
	GetAllExpenses() ([]Domain.Expense, error)
	GetExpenseByID(id string) (*Domain.Expense, error)
	VerifyExpense(id string) error
}

type expenseUsecase struct {
	repo Repositories.ExpenseRepository
}

func NewExpenseUsecase(repo Repositories.ExpenseRepository) ExpenseUsecase {
	return &expenseUsecase{repo: repo}
}

func (u *expenseUsecase) CreateExpense(input *Domain.Expense) (*Domain.Expense, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	input.Status = "pending"
	input.CreatedAt = time.Now().UTC()
	if err := u.repo.Create(input); err != nil {
		return nil, err
	}
	return input, nil
}

func (u *expenseUsecase) GetAllExpenses() ([]Domain.Expense, error) {
	return u.repo.GetAll()
}

func (u *expenseUsecase) GetExpenseByID(id string) (*Domain.Expense, error) {
	return u.repo.GetByID(id)
}

func (u *expenseUsecase) VerifyExpense(id string) error {
	e, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	e.Status = "verified"
	return u.repo.Update(id, e)
}
