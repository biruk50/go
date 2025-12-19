package Usecases

import (
	"errors"
	"testing"

	"FMS/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mock cash repo
type mockCashRepo struct {
	store map[string]*Domain.CashRequest
}

func newMockCashRepo() *mockCashRepo {
	return &mockCashRepo{store: make(map[string]*Domain.CashRequest)}
}
func (m *mockCashRepo) Create(t *Domain.CashRequest) error {
	if t == nil {
		return errors.New("nil")
	}
	m.store[t.ID.Hex()] = t
	return nil
}
func (m *mockCashRepo) GetAll() ([]Domain.CashRequest, error) {
	res := make([]Domain.CashRequest, 0, len(m.store))
	for _, v := range m.store {
		res = append(res, *v)
	}
	return res, nil
}
func (m *mockCashRepo) GetByID(id string) (*Domain.CashRequest, error) {
	if v, ok := m.store[id]; ok {
		return v, nil
	}
	return nil, errors.New("not found")
}
func (m *mockCashRepo) Update(id string, t *Domain.CashRequest) error {
	if _, ok := m.store[id]; !ok {
		return errors.New("not found")
	}
	m.store[id] = t
	return nil
}
func (m *mockCashRepo) Delete(id string) error {
	if _, ok := m.store[id]; !ok {
		return errors.New("not found")
	}
	delete(m.store, id)
	return nil
}

func TestCashUsecase_CreateApproveDisburse(t *testing.T) {
	mock := newMockCashRepo()
	uc := NewCashRequestUsecase(mock)

	c := &Domain.CashRequest{ID: primitive.NewObjectID(), Title: "Req1", Amount: 500}
	created, err := uc.CreateCashRequest(c)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Status != "pending" {
		t.Fatalf("expected pending")
	}

	if err := uc.ApproveCashRequest(created.ID.Hex()); err != nil {
		t.Fatalf("approve failed: %v", err)
	}
	cr, _ := uc.GetCashRequestByID(created.ID.Hex())
	if cr.Status != "approved" {
		t.Fatalf("expected approved")
	}

	if err := uc.DisburseCashRequest(created.ID.Hex()); err != nil {
		t.Fatalf("disburse failed: %v", err)
	}
	cr2, _ := uc.GetCashRequestByID(created.ID.Hex())
	if cr2.Status != "disbursed" {
		t.Fatalf("expected disbursed")
	}
}
