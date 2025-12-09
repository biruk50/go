package mocks

import (
	"github.com/stretchr/testify/mock"
	"task_manager_testify/Domain"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(t *Domain.Task) error {
	args := m.Called(t) 
	return args.Error(0)
}

func (m *MockTaskRepository) GetAll() ([]Domain.Task, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]Domain.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) GetByID(id string) (*Domain.Task, error) {
	args := m.Called(id)
	if obj := args.Get(0); obj != nil {
		return obj.(*Domain.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) Update(id string, t *Domain.Task) error {
	args := m.Called(id, t)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
