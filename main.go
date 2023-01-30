package main

import (
	"net/http"

	"github.com/makcik45/jwt-go/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/register", controllers.Register)

	// fmt.Println("Server jalan di: http://localhost:3000")
	http.ListenAndServe(":8080", nil)
}
