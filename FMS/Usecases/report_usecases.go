package Usecases

import (
	"FMS/Domain"
	"FMS/Repositories"
)

type ReportUsecase interface {
	GetOverview() (map[string]interface{}, error)
	GetBudgetReport() ([]Domain.Budget, error)
	GetCashRequestReport() ([]Domain.CashRequest, error)
	GetExpenseReport() ([]Domain.Expense, error)
}

type reportUsecase struct {
	budgetRepo      Repositories.BudgetRepository
	cashRequestRepo Repositories.CashRequestRepository
	expenseRepo     Repositories.ExpenseRepository
}

func NewReportUsecase(b Repositories.BudgetRepository, c Repositories.CashRequestRepository, e Repositories.ExpenseRepository) ReportUsecase {
	return &reportUsecase{budgetRepo: b, cashRequestRepo: c, expenseRepo: e}
}

func (u *reportUsecase) GetOverview() (map[string]interface{}, error) {
	budgets, _ := u.budgetRepo.GetAll()
	cashReqs, _ := u.cashRequestRepo.GetAll()
	expenses, _ := u.expenseRepo.GetAll()

	overview := map[string]interface{}{
		"budgets_count":       len(budgets),
		"cash_requests_count": len(cashReqs),
		"expenses_count":      len(expenses),
	}
	return overview, nil
}

func (u *reportUsecase) GetBudgetReport() ([]Domain.Budget, error) {
	return u.budgetRepo.GetAll()
}

func (u *reportUsecase) GetCashRequestReport() ([]Domain.CashRequest, error) {
	return u.cashRequestRepo.GetAll()
}

func (u *reportUsecase) GetExpenseReport() ([]Domain.Expense, error) {
	return u.expenseRepo.GetAll()
}
