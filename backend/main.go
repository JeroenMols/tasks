package main

import (
	"backend/db"
	"backend/net"
	"backend/routes"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	database := db.InMemoryDatabase()

	users := routes.Users{Database: database}
	todoLists := routes.TodoLists{Database: database}
	todos := routes.Todos{Database: database}

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)

	mux.HandleFunc("POST /todolists", todoLists.Create)
	mux.HandleFunc("GET /todolists/{list_id}", todoLists.Get)

	mux.HandleFunc("POST /todos", todos.Create)
	mux.HandleFunc("PUT /todos/{todo_id}", todos.Update)

	// Debug route
	debug := routes.Debug{Database: database}
	mux.HandleFunc("GET /debug", debug.Debug)

	handler := net.CorsMiddleware(mux, "*")

	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		fmt.Println(err.Error())
	}
}
