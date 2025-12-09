package routers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"task_manager_testify/Delivery/routers"
	"task_manager_testify/Domain"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Tests/mocks"
	"task_manager_testify/Usecases"
)

type RouterTestSuite struct {
	suite.Suite
	router     *gin.Engine
	userRepo   *mocks.MockUserRepository
	taskRepo   *mocks.MockTaskRepository
	jwtService Infrastructure.JWTService
}

func (s *RouterTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	s.userRepo = &mocks.MockUserRepository{}
	s.taskRepo = &mocks.MockTaskRepository{}
	s.jwtService = Infrastructure.NewJWTService("router-suite-secret")

	pw := Infrastructure.NewPasswordService()
	userUC := Usecases.NewUserUsecase(s.userRepo, pw, s.jwtService)
	taskUC := Usecases.NewTaskUsecase(s.taskRepo)

	s.router = routers.SetupRouter(userUC, taskUC, s.jwtService)
}

func (s *RouterTestSuite) TestInvalidJSONRegister() {
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(400, w.Code)
}

func (s *RouterTestSuite) TestMissingFieldsLogin() {
	body := map[string]string{"username": "alice"} // missing password
	b, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(400, w.Code)
}

func (s *RouterTestSuite) TestProtectedRoute_NoToken() {
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(401, w.Code)
}

func (s *RouterTestSuite) TestAdminRoute_NonAdminDenied() {
	pw := Infrastructure.NewPasswordService()
	hash, _ := pw.Hash("pass")

	s.userRepo.On("FindByUsername", "bob").Return(&Domain.User{
		Username:     "bob",
		PasswordHash: hash,
		Role:         "user",
	}, nil)

	token, _ := s.jwtService.Generate("1", "bob", "user", time.Hour)

	req := httptest.NewRequest("POST", "/promote?username=alice", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(403, w.Code)
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
