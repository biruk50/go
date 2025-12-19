package Usecases

import (
	"FMS/Domain"
	"FMS/Repositories"
	"errors"
	"time"
)

// BudgetUsecase defines business operations for budgets
type BudgetUsecase interface {
	CreateBudget(input *Domain.Budget) (*Domain.Budget, error)
	GetAllBudgets() ([]Domain.Budget, error)
	GetBudgetByID(id string) (*Domain.Budget, error)
	GetBudgetSummary(id string) (map[string]interface{}, error)
	UpdateBudget(id string, input *Domain.Budget) error
	ApproveBudget(id string) error
	RejectBudget(id string) error
}

type budgetUsecase struct {
	budgetRepo Repositories.BudgetRepository
}

func NewBudgetUsecase(repo Repositories.BudgetRepository) BudgetUsecase {
	return &budgetUsecase{budgetRepo: repo}
}

func (u *budgetUsecase) CreateBudget(input *Domain.Budget) (*Domain.Budget, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.DueDate.IsZero() {
		input.DueDate = time.Now().Add(24 * time.Hour)
	}

	// initialize status and remaining amount
	input.Status = "pending"
	input.Remaining = input.Amount

	if err := u.budgetRepo.Create(input); err != nil {
		return nil, err
	}

	return input, nil
}

func (u *budgetUsecase) GetAllBudgets() ([]Domain.Budget, error) {
	return u.budgetRepo.GetAll()
}

func (u *budgetUsecase) GetBudgetByID(id string) (*Domain.Budget, error) {
	return u.budgetRepo.GetByID(id)
}

func (u *budgetUsecase) GetBudgetSummary(id string) (map[string]interface{}, error) {
	b, err := u.budgetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	summary := map[string]interface{}{
		"id":        b.ID.Hex(),
		"title":     b.Title,
		"amount":    b.Amount,
		"remaining": b.Remaining,
		"status":    b.Status,
	}
	return summary, nil
}

func (u *budgetUsecase) UpdateBudget(id string, input *Domain.Budget) error {
	if input.Title == "" {
		return errors.New("title is required")
	}
	return u.budgetRepo.Update(id, input)
}

func (u *budgetUsecase) ApproveBudget(id string) error {
	b, err := u.budgetRepo.GetByID(id)
	if err != nil {
		return err
	}
	b.Status = "approved"
	b.Remaining = b.Amount
	return u.budgetRepo.Update(id, b)
}

func (u *budgetUsecase) RejectBudget(id string) error {
	b, err := u.budgetRepo.GetByID(id)
	if err != nil {
		return err
	}
	b.Status = "rejected"
	return u.budgetRepo.Update(id, b)
}
