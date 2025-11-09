package router

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", controllers.Home)

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", controllers.GetAllTasks)
		taskRoutes.GET("/:id", controllers.GetTask)
		taskRoutes.POST("", controllers.CreateTask)
		taskRoutes.PUT("/:id", controllers.UpdateTask)
		taskRoutes.DELETE("/:id", controllers.DeleteTask)
	}

	return router
}
