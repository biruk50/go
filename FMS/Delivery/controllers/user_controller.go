package controllers

import (
	"FMS/Infrastructure"
	"FMS/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUC Usecases.UserUsecase
	JWT    Infrastructure.JWTService
}

func NewUserController(u Usecases.UserUsecase, jwt Infrastructure.JWTService) *UserController {
	return &UserController{UserUC: u, JWT: jwt}
}

func (uc *UserController) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the FMS API"})
}

func (uc *UserController) Register(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uc.UserUC.Register(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"username": user.Username, "role": user.Role, "created_at": user.CreatedAt})
}

func (uc *UserController) Login(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := uc.UserUC.Login(payload.Username, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expires_in": 24 * 3600})
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	// Not implemented fully - passthrough to UserUC if available
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (uc *UserController) GetMyProfile(c *gin.Context) {
	uname, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"username": uname})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role required"})
		return
	}
	if err := uc.UserUC.Promote(payload.Role); err != nil {
		// Promote is being reused; adjust in real implementation
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role updated", "id": id})
}
