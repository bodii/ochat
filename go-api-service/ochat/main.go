package main

import (
	"net/http"
	"ochat/bootstrap"
	"ochat/controllers"
)

// main func
func main() {
	bootstrap.Init()

	http.HandleFunc("/user/login", controllers.Login)
	http.HandleFunc("/user/signup", controllers.Register)

	http.ListenAndServe(":8080", nil)
}
