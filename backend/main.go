package main

import (
	"backend/db"
	"backend/net"
	"backend/routes"
	"backend/util"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	database := db.InMemoryDatabase()

	users := routes.Users{
		Database:     database,
		GenerateUuid: util.GenerateRandomUuid,
	}
	todoLists := routes.TodoLists{
		Database:     database,
		GenerateUuid: util.GenerateRandomUuid,
	}
	todos := routes.Todos{
		Database:     database,
		GenerateUuid: util.GenerateRandomUuid,
		CurrentTime:  util.GetCurrentTime,
	}

	mux.HandleFunc("POST /users/register", users.Register)
	mux.HandleFunc("POST /users/login", users.Login)

	mux.HandleFunc("POST /todolists", todoLists.Create)
	mux.HandleFunc("GET /todolists/{list_id}", todoLists.Get)

	mux.HandleFunc("POST /todos", todos.Create)
	mux.HandleFunc("PUT /todos", todos.Update)

	// Debug route
	debug := routes.Debug{
		Database: database,
	}
	mux.HandleFunc("GET /debug", debug.Debug)

	handler := net.CorsMiddleware(mux, "*")

	err := http.ListenAndServe("localhost:8080", handler)
	if err != nil {
		fmt.Println(err.Error())
	}
}
