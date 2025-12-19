package Usecases

import (
	"errors"
	"time"

	"task_manager_clean/Domain"
	"task_manager_clean/Repositories"
)

type TaskUsecase interface {
	CreateTask(input *Domain.Task) (*Domain.Task, error)
	GetTasks() ([]Domain.Task, error)
	GetTaskByID(id string) (*Domain.Task, error)
	UpdateTask(id string, input *Domain.Task) error
	DeleteTask(id string) error
}

type taskUsecase struct {
	taskRepo Repositories.TaskRepository
}

func NewTaskUsecase(repo Repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{taskRepo: repo}
}

func (u *taskUsecase) CreateTask(input *Domain.Task) (*Domain.Task, error) {
	if input.Title == "" {
		return nil, errors.New("title is required")
	}

	if input.DueDate.IsZero() {
		input.DueDate = time.Now().Add(24 * time.Hour)
	}

	if err := u.taskRepo.Create(input); err != nil {
		return nil, err
	}

	return input, nil
}

func (u *taskUsecase) GetTasks() ([]Domain.Task, error) {
	return u.taskRepo.GetAll()
}

func (u *taskUsecase) GetTaskByID(id string) (*Domain.Task, error) {
	return u.taskRepo.GetByID(id)
}

func (u *taskUsecase) UpdateTask(id string, input *Domain.Task) error {
	if input.Title == "" {
		return errors.New("title is required")
	}
	return u.taskRepo.Update(id, input)
}

func (u *taskUsecase) DeleteTask(id string) error {
	return u.taskRepo.Delete(id)
}

