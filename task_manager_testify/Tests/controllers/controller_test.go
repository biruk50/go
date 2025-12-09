package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"task_manager_testify/Delivery/routers"
	"task_manager_testify/Domain"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Tests/mocks"
	"task_manager_testify/Usecases"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterAndLoginRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// mocks
	userRepo := &mocks.MockUserRepository{}
	pw := Infrastructure.NewPasswordService()
	j := Infrastructure.NewJWTService("test-secret")
	userUC := Usecases.NewUserUsecase(userRepo, pw, j)
	taskUC := Usecases.NewTaskUsecase(nil) // not used here

	router := routers.SetupRouter(userUC, taskUC, j) // adjust signature

	// Register: ensure FindByUsername returns not-found, then Count/Create
	userRepo.On("FindByUsername", "alice").Return(nil, errors.New("not found")).Once()
	userRepo.On("Count").Return(int64(0), nil).Once()
	userRepo.On("Create", mock.Anything).Return(nil).Once()

	payload := map[string]string{"username": "alice", "password": "pwd"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	// Login: find user
	hashed, _ := pw.Hash("pwd")
	userRepo.On("FindByUsername", "alice").Return(&Domain.User{Username: "alice", PasswordHash: hashed, Role: "admin"}, nil).Once()

	payloadLogin := map[string]string{"username": "alice", "password": "pwd"}
	b2, _ := json.Marshal(payloadLogin)
	req2 := httptest.NewRequest("POST", "/login", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	taskRepo := &mocks.MockTaskRepository{}
	userRepo := &mocks.MockUserRepository{}
	pw := Infrastructure.NewPasswordService()
	j := Infrastructure.NewJWTService("bad-json")

	userUC := Usecases.NewUserUsecase(userRepo, pw, j)
	taskUC := Usecases.NewTaskUsecase(taskRepo)
	r := routers.SetupRouter(userUC, taskUC, j)

	token, _ := j.Generate("1", "bob", "admin", time.Hour)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("{bad json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}
