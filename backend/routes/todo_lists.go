package routes

import (
	"backend/db"
	"backend/models"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
	"time"
)

type TodoLists struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
}

func (t *TodoLists) Create(w http.ResponseWriter, r *http.Request) {
	_, err := net.ParseBody[listCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}
	if _, err := t.Database.Authorize(r.Header.Get("Authorization")); err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	listId := t.Database.CreateTodoList(t.GenerateUuid())
	fmt.Printf("Creating new todo list %s\n", listId)

	net.Success(w, listCreateResponse{TodoListId: listId})
}

func (t *TodoLists) Get(w http.ResponseWriter, r *http.Request) {
	if _, err := t.Database.Authorize(r.Header.Get("Authorization")); err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	listId := r.PathValue("list_id")

	todos, err := t.Database.GetTodos(listId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	formattedTodos := []models.TodoItem{}
	for _, todo := range *todos {
		formattedTodos = append(formattedTodos, models.TodoItem{
			Id:          todo.Id,
			CreatedBy:   t.Database.Users[todo.User],
			Description: todo.Description,
			Status:      todo.Status,
			UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
		})
	}

	net.Success(w, listGetResponse{Todos: formattedTodos})
}

type listCreateRequest struct {
}

type listCreateResponse struct {
	TodoListId string `json:"todo_list_id"`
}

// TODO include the todo list id here
type listGetResponse struct {
	Todos []models.TodoItem `json:"todos"`
}
