package main

import (
  "Task_Management/router"
)

func main() {
	router := router.SetupRouter()

	router.Run(":8080")
  
}

