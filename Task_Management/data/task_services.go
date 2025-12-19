package data

import (
	"Task_Management/models"
	"errors"
)

var tasks []models.Task

// GetAllTasks returns all tasks.
func GetAllTasks() []models.Task {
	return tasks
}

// GetTaskByID returns a single task by ID.
func GetTaskByID(id string) (*models.Task, error) {
	for _, task := range tasks {
		if task.Id == id {
			return &task, nil
		}
	}
	return nil, errors.New("task not found")
}

func AddTask(task models.Task) {
	tasks = append(tasks, task)
}

func UpdateTask(id string, updated models.Task) error {
	for i, t := range tasks {
		if t.Id == id {
			tasks[i] = updated
			return nil
		}
	}
	return errors.New("task not found")
}

func DeleteTask(id string) error {
	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
