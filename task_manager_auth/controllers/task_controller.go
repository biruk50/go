package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// Home route.
func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the Task Manager API"})
}

// Get all tasks.
func GetAllTasks(ctx *gin.Context) {
	tasks,err := data.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks })
}

// Get task by ID.
func GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

// Create a new task.
func CreateTask(ctx *gin.Context) {
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newTask.ID.IsZero() {
		newTask.ID = primitive.NewObjectID()
	}

	if newTask.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "title is required"})
		return
	}

	// Default due date
	if newTask.DueDate.IsZero() {
		newTask.DueDate = time.Now().Add(24 * time.Hour)
	}

	if err := data.AddTask(newTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task created",
		"task":    newTask,
	})
}

// Update an existing task.
func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updated models.Task
	if err := ctx.ShouldBindJSON(&updated); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.UpdateTask(id, updated); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// Delete a task.
func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := data.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
