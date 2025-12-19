package main

import (
	"log"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	if err := data.InitMongo(); err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}
	defer data.CloseMongo()

	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
