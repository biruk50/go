package controllers

import (
	"net/http"
	"task_manager_clean/Domain"
	"task_manager_clean/Usecases"
	"task_manager_clean/Infrastructure"

	"github.com/gin-gonic/gin"
)

// Controllers depend only on usecases (interfaces)
type Controller struct {
	UserUC Usecases.UserUsecase
	TaskUC Usecases.TaskUsecase
	JWT    Infrastructure.JWTService
}

// NewController factory - note: we pass JWT service for middleware token examples
func NewController(u Usecases.UserUsecase, t Usecases.TaskUsecase, jwtSvc Infrastructure.JWTService) *Controller {
	return &Controller{UserUC: u, TaskUC: t, JWT: jwtSvc}
}

func (ctr *Controller) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Task Manager API"})
}

func (ctr *Controller) Register(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ctr.UserUC.Register(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "role": user.Role, "created_at": user.CreatedAt})
}


func (ctr *Controller) Login(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := ctr.UserUC.Login(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expires_in": 24 * 3600})
}

// Promote (admin)
func (ctr *Controller) Promote(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	if err := ctr.UserUC.Promote(payload.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "promoted"})
}

// Tasks handlers
func (ctr *Controller) GetAllTasks(c *gin.Context) {
	tasks, err := ctr.TaskUC.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (ctr *Controller) GetTask(c *gin.Context) {
	id := c.Param("id")
	t, err := ctr.TaskUC.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": t})
}

func (ctr *Controller) CreateTask(c *gin.Context) {
	var payload Domain.Task
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.DueDate = payload.DueDate.UTC()
	created, err := ctr.TaskUC.CreateTask(&payload)
	if err != nil {
		if _, ok := err.(interface{ Error() string }); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": created})
}

func (ctr *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var payload Domain.Task
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctr.TaskUC.UpdateTask(id, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (ctr *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := ctr.TaskUC.DeleteTask(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

