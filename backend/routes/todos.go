package routes

import (
	"backend/db"
	"backend/models"
	"backend/net"
	"backend/util"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type Todos struct {
	Database     db.Database
	GenerateUuid util.GenerateUuid
	CurrentTime  util.CurrentTime
}

func (t *Todos) Create(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoCreateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	accountNumber, err := t.Database.Authorize(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	if !regexp.MustCompile(todoDescriptionRegex).MatchString(body.Description) {
		net.HaltBadRequest(w, "invalid description")
		return
	}

	todos, err := t.Database.GetTodos(body.ListId)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	item := models.TodoDatabaseItem{
		Id:          t.GenerateUuid(),
		Description: body.Description,
		Status:      "todo",
		User:        *accountNumber,
		UpdatedAt:   t.CurrentTime(),
	}
	fmt.Printf("Creating new todo %s\n", item.Description)
	t.Database.TodoLists[body.ListId] = append(*todos, item)

	net.Success(w, models.TodoItem{
		Id:          item.Id,
		CreatedBy:   t.Database.Users[*accountNumber],
		Description: item.Description,
		Status:      item.Status,
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	})
}

func (t *Todos) Update(w http.ResponseWriter, r *http.Request) {
	body, err := net.ParseBody[todoUpdateRequest](r)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	todo_id := r.PathValue("todo_id")

	accountNumber, err := t.Database.Authorize(r.Header.Get("Authorization"))
	if err != nil {
		net.HaltUnauthorized(w, err.Error())
		return
	}

	item, err := t.Database.GetTodo(todo_id)
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	err = item.ChangeStatus(body.Status, t.CurrentTime())
	if err != nil {
		net.HaltBadRequest(w, err.Error())
		return
	}

	net.Success(w, item.ToTodoItem(t.Database.Users[*accountNumber]))
}

type todoCreateRequest struct {
	Description string `json:"description" validate:"required"`
	ListId      string `json:"todo_list_id" validate:"required"`
}

type todoUpdateRequest struct {
	Status string `json:"status" validate:"required"`
}
