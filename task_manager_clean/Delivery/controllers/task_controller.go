 package controllers

import (
	"net/http"
	"task_manager_clean/data"
	"task_manager_clean/models"
	"task_manager_clean/middleware"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// Register request payload
type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login request
type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register endpoint
func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := data.CreateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "role": user.Role, "created_at": user.CreatedAt})
}

// Login endpoint -> returns JWT
func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := data.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	
	// create token
	secret := data.GetJWTSecret()
	claims := middleware.Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.Hex(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed, "expires_in": 24 * 3600})
}

// Promote endpoint (admin-only). Request body expects {"username":"otheruser"}
func Promote(c *gin.Context) {
	var payload struct{ Username string `json:"username"` }
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	if err := data.PromoteUser(payload.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}


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
