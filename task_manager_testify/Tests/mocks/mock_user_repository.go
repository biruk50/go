package mocks

import (
	"github.com/stretchr/testify/mock"
	"task_manager_testify/Domain"
)

// MockUserRepository is a testify mock for user repo interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *Domain.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*Domain.User, error) {
	args := m.Called(username)
	if obj := args.Get(0); obj != nil {
		return obj.(*Domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Promote(username string) error {
	args := m.Called(username)
	return args.Error(0)
}
