package main

import (
	"fmt"
	"user-service/api"
	"user-service/db"
	"user-service/error_tracer"
)

func main() {
	initErrorTracer()
	initDB()
	defer db.Client.Close()

	r := api.NewGinEngine()
	fmt.Println("User service is running on port :8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func initDB() {
	db.Client = &db.UserRepository{}
	if err := db.Client.Init(); err != nil {
		panic(err)
	}
}

func initErrorTracer() {
	error_tracer.Client = &error_tracer.ConsoleLog{}
}
