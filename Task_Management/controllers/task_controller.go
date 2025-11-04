package controllers

import (
	"net/http"
	"Task_Management/data"
	"Task_Management/models"
	"time"

	"github.com/gin-gonic/gin"
)

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to the Task Manager API"})
}

func GetAllTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"tasks": data.GetAllTasks()})
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

func CreateTask(ctx *gin.Context) {
	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if newTask.Id == "" || newTask.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "id and title are required"})
		return
	}

	// Default due date
	if newTask.DueDate.IsZero() {
		newTask.DueDate = time.Now().Add(24 * time.Hour)
	}

	data.AddTask(newTask)
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

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

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := data.DeleteTask(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
