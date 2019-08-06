package main

import (
	"fmt"
	"password-service/api"
	"password-service/db"
	"password-service/error_tracer"
)

func main() {
	initErrorTracer()
	initDB()
	defer db.Client.Close()

	r := api.NewGinEngine()
	fmt.Println("Password service is running on port :8081")
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}

func initDB() {
	db.Client = &db.PostgresClient{}
	if err := db.Client.Init(); err != nil {
		panic(err)
	}
}

func initErrorTracer() {
	error_tracer.Client = &error_tracer.ConsoleLog{}
}
