package router

import (
	"task_manager_auth/controllers"
	"task_manager_auth/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public
	router.GET("/", controllers.Home)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Protected: any authenticated user
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/tasks", controllers.GetAllTasks)
		auth.GET("/tasks/:id", controllers.GetTask)
	}

	// Admin-only
	admin := router.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.POST("/tasks", controllers.CreateTask)
		admin.PUT("/tasks/:id", controllers.UpdateTask)
		admin.DELETE("/tasks/:id", controllers.DeleteTask)

		admin.POST("/promote", controllers.Promote)
	}

	return router

}
