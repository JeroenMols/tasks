package routes

import (
	"backend/db"
	"backend/net"
	"net/http"
	"regexp"
)

type Todos struct {
	database db.Database
}

func CreateTodos(database db.Database) Todos {
	return Todos{database: database}
}

func (t *Todos) Create(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	accountNumber, err := t.database.Authorize(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	if !regexp.MustCompile(todoDescriptionRegex).MatchString(body.Description) {
		net.HaltBadRequest(w, "invalid description")
		return
	}

	_, err = t.database.GetTodos(body.ListId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	item := t.database.CreateTodo(body.ListId, body.Description, *accountNumber)

	net.Success(w, item.ToTodoItem(t.database.Users[*accountNumber].Name))
}

func (t *Todos) Update(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoUpdateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	todo_id := r.PathValue("todo_id")

	accountNumber, err := t.database.Authorize(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	item, err := t.database.GetTodo(todo_id)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	err = item.ChangeStatus(body.Status)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}
	t.database.UpdateTodo(item)

	net.Success(w, item.ToTodoItem(t.database.Users[*accountNumber].Name))
}

type todoCreateRequest struct {
	Description string `json:"description" validate:"required"`
	ListId      string `json:"todo_list_id" validate:"required"`
}

type todoUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}
