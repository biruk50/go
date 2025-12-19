package Usecases

import (
	"FMS/Domain"
	"FMS/Repositories"
	"errors"
	"time"
)

type CashRequestUsecase interface {
	CreateCashRequest(input *Domain.CashRequest) (*Domain.CashRequest, error)
	GetAllCashRequests() ([]Domain.CashRequest, error)
	GetCashRequestByID(id string) (*Domain.CashRequest, error)
	ApproveCashRequest(id string) error
	RejectCashRequest(id string) error
	DisburseCashRequest(id string) error
}

type cashRequestUsecase struct {
	repo Repositories.CashRequestRepository
}

func NewCashRequestUsecase(repo Repositories.CashRequestRepository) CashRequestUsecase {
	return &cashRequestUsecase{repo: repo}
}

func (u *cashRequestUsecase) CreateCashRequest(input *Domain.CashRequest) (*Domain.CashRequest, error) {
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

func (u *cashRequestUsecase) GetAllCashRequests() ([]Domain.CashRequest, error) {
	return u.repo.GetAll()
}

func (u *cashRequestUsecase) GetCashRequestByID(id string) (*Domain.CashRequest, error) {
	return u.repo.GetByID(id)
}

func (u *cashRequestUsecase) ApproveCashRequest(id string) error {
	r, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	r.Status = "approved"
	return u.repo.Update(id, r)
}

func (u *cashRequestUsecase) RejectCashRequest(id string) error {
	r, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	r.Status = "rejected"
	return u.repo.Update(id, r)
}

func (u *cashRequestUsecase) DisburseCashRequest(id string) error {
	r, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if r.Status != "approved" {
		return errors.New("only approved requests can be disbursed")
	}
	r.Status = "disbursed"
	return u.repo.Update(id, r)
}
