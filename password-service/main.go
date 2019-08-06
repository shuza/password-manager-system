package main

import (
	"fmt"
	"password-service/api"
)

func main() {
	r := api.NewGinEngine()
	fmt.Println("Password service is running on port :8081")
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
