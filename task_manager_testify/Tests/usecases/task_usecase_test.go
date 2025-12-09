package usecases

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"task_manager_testify/Domain"
	"task_manager_testify/Tests/mocks"
	"task_manager_testify/Usecases"
)

type TaskUsecaseSuite struct {
	suite.Suite
	taskRepo *mocks.MockTaskRepository
	usecase  Usecases.TaskUsecase
}

func (s *TaskUsecaseSuite) SetupSuite() {
	s.taskRepo = &mocks.MockTaskRepository{}
	s.usecase = Usecases.NewTaskUsecase(s.taskRepo)
}

func (s *TaskUsecaseSuite) TestCreateTask_Success() {
	input := Domain.Task{Title: "Buy milk"}
	s.taskRepo.On("Create", mock.Anything).Return(nil)

	res, err := s.usecase.CreateTask(&input)

	s.NoError(err)
	s.Equal("Buy milk", res.Title)
}

func (s *TaskUsecaseSuite) TestCreateTask_MissingTitle() {
	_, err := s.usecase.CreateTask(&Domain.Task{})
	s.Error(err)
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}
