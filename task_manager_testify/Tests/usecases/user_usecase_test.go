package usecases

import (
	"errors"
	"testing"

	"task_manager_testify/Domain"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Tests/mocks"
	"task_manager_testify/Usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Login_Promote(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	pwSvc := Infrastructure.NewPasswordService()
	jwtSvc := Infrastructure.NewJWTService("jwt-secret")

	uc := Usecases.NewUserUsecase(mockRepo, pwSvc, jwtSvc)

	// Register success: make FindByUsername return not-found, repo.Count returns 0 -> first user becomes admin
	mockRepo.On("FindByUsername", "alice").Return(nil, errors.New("not found"))
	mockRepo.On("Count").Return(int64(0), nil)
	mockRepo.On("Create", mock.AnythingOfType("*Domain.User")).Return(nil)

	u, err := uc.Register("alice", "pwd")
	assert.NoError(t, err)
	assert.Equal(t, "alice", u.Username)
	// Since Create was stubbed, Register will set role "admin"
	mockRepo.AssertCalled(t, "Create", mock.Anything)

	// Register duplicate
	mockRepo2 := &mocks.MockUserRepository{}
	mockRepo2.On("FindByUsername", "bob").Return(&Domain.User{Username: "bob"}, nil)
	uc2 := Usecases.NewUserUsecase(mockRepo2, pwSvc, jwtSvc)
	_, err = uc2.Register("bob", "pwd")
	assert.Error(t, err)

	// Login success
	mockRepo3 := &mocks.MockUserRepository{}
	hashed, _ := pwSvc.Hash("pwd")
	mockRepo3.On("FindByUsername", "carl").Return(&Domain.User{ID: Domain.User{}.ID, Username: "carl", PasswordHash: hashed, Role: "user"}, nil)
	uc3 := Usecases.NewUserUsecase(mockRepo3, pwSvc, jwtSvc)
	token, err := uc3.Login("carl", "pwd")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Login failure bad pw
	mockRepo3.AssertCalled(t, "FindByUsername", "carl")
}
