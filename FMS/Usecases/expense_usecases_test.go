package Usecases

import (
	"errors"
	"testing"

	"FMS/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mock expense repo
type mockExpenseRepo struct{ store map[string]*Domain.Expense }

func newMockExpenseRepo() *mockExpenseRepo {
	return &mockExpenseRepo{store: make(map[string]*Domain.Expense)}
}
func (m *mockExpenseRepo) Create(t *Domain.Expense) error {
	if t == nil {
		return errors.New("nil")
	}
	m.store[t.ID.Hex()] = t
	return nil
}
func (m *mockExpenseRepo) GetAll() ([]Domain.Expense, error) {
	res := make([]Domain.Expense, 0, len(m.store))
	for _, v := range m.store {
		res = append(res, *v)
	}
	return res, nil
}
func (m *mockExpenseRepo) GetByID(id string) (*Domain.Expense, error) {
	if v, ok := m.store[id]; ok {
		return v, nil
	}
	return nil, errors.New("not found")
}
func (m *mockExpenseRepo) Update(id string, t *Domain.Expense) error {
	if _, ok := m.store[id]; !ok {
		return errors.New("not found")
	}
	m.store[id] = t
	return nil
}
func (m *mockExpenseRepo) Delete(id string) error {
	if _, ok := m.store[id]; !ok {
		return errors.New("not found")
	}
	delete(m.store, id)
	return nil
}

func TestExpenseUsecase_CreateAttachVerify(t *testing.T) {
	mock := newMockExpenseRepo()
	uc := NewExpenseUsecase(mock)

	e := &Domain.Expense{ID: primitive.NewObjectID(), Title: "Lunch", Amount: 20}
	created, err := uc.CreateExpense(e)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Status != "pending" {
		t.Fatalf("expected pending")
	}

	if err := uc.AttachReceipt(created.ID.Hex(), "http://example.com/rec.jpg"); err != nil {
		t.Fatalf("attach failed: %v", err)
	}
	e2, _ := uc.GetExpenseByID(created.ID.Hex())
	if e2.ReceiptURL != "http://example.com/rec.jpg" {
		t.Fatalf("receipt not attached")
	}

	if err := uc.VerifyExpense(created.ID.Hex()); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	e3, _ := uc.GetExpenseByID(created.ID.Hex())
	if e3.Status != "verified" {
		t.Fatalf("expected verified")
	}
}
