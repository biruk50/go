package Usecases

import (
	"testing"

	"FMS/Domain"
)

type mockBudgetRepoSimple struct{ store []Domain.Budget }

func (m *mockBudgetRepoSimple) Create(t *Domain.Budget) error             { return nil }
func (m *mockBudgetRepoSimple) GetAll() ([]Domain.Budget, error)          { return m.store, nil }
func (m *mockBudgetRepoSimple) GetByID(id string) (*Domain.Budget, error) { return nil, nil }
func (m *mockBudgetRepoSimple) Update(id string, t *Domain.Budget) error  { return nil }
func (m *mockBudgetRepoSimple) Delete(id string) error                    { return nil }

type mockCashRepoSimple struct{ store []Domain.CashRequest }

func (m *mockCashRepoSimple) Create(t *Domain.CashRequest) error             { return nil }
func (m *mockCashRepoSimple) GetAll() ([]Domain.CashRequest, error)          { return m.store, nil }
func (m *mockCashRepoSimple) GetByID(id string) (*Domain.CashRequest, error) { return nil, nil }
func (m *mockCashRepoSimple) Update(id string, t *Domain.CashRequest) error  { return nil }
func (m *mockCashRepoSimple) Delete(id string) error                         { return nil }

type mockExpenseRepoSimple struct{ store []Domain.Expense }

func (m *mockExpenseRepoSimple) Create(t *Domain.Expense) error             { return nil }
func (m *mockExpenseRepoSimple) GetAll() ([]Domain.Expense, error)          { return m.store, nil }
func (m *mockExpenseRepoSimple) GetByID(id string) (*Domain.Expense, error) { return nil, nil }
func (m *mockExpenseRepoSimple) Update(id string, t *Domain.Expense) error  { return nil }
func (m *mockExpenseRepoSimple) Delete(id string) error                     { return nil }

func TestReportUsecase_GetOverview(t *testing.T) {
	b := mockBudgetRepoSimple{store: []Domain.Budget{{}}}
	c := mockCashRepoSimple{store: []Domain.CashRequest{{}, {}}}
	e := mockExpenseRepoSimple{store: []Domain.Expense{{}}}

	ru := NewReportUsecase(&b, &c, &e)
	o, err := ru.GetOverview()
	if err != nil {
		t.Fatalf("overview failed: %v", err)
	}
	if o["budgets_count"].(int) != 1 {
		t.Fatalf("budgets_count expected 1")
	}
	if o["cash_requests_count"].(int) != 2 {
		t.Fatalf("cash_requests_count expected 2")
	}
	if o["expenses_count"].(int) != 1 {
		t.Fatalf("expenses_count expected 1")
	}
}
