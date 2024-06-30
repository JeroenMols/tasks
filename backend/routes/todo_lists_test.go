package routes

import (
	"backend/db"
	"backend/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type createListTestCase struct {
	description   string
	accessToken   string
	body          string
	responseCode  int
	responseBody  string
	databaseLists map[string][]models.TodoDatabaseItem
}

func TestTodoLists_Create(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-222222222222"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-333333333333"

	tests := []createListTestCase{
		{
			description:   "Invalid body",
			accessToken:   validAccessToken,
			body:          `{"invalid":"body"}`,
			responseCode:  http.StatusBadRequest,
			responseBody:  `{"error":"invalid body"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description:   "Access token does not exist",
			accessToken:   nonExistingAccessToken,
			body:          `{}`,
			responseCode:  http.StatusUnauthorized,
			responseBody:  `{"error":"account not found"}`,
			databaseLists: make(map[string][]models.TodoDatabaseItem),
		},
		{
			description:   "Create new todo list",
			accessToken:   validAccessToken,
			body:          `{}`,
			responseCode:  http.StatusOK,
			responseBody:  `{"todo_list_id":"static_uuid"}`,
			databaseLists: map[string][]models.TodoDatabaseItem{"static_uuid": {}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.AccessTokens[validAccessToken] = existingAccount

			todoList := TodoLists{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
			}

			request := httptest.NewRequest(http.MethodPost, "/todolists", strings.NewReader(tt.body))
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todoList.Create(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
			assert.Equal(t, tt.databaseLists, database.TodoLists)
		})
	}
}

type getListTestCase struct {
	description  string
	accessToken  string
	todoListId   string
	responseCode int
	responseBody string
}

func TestTodoLists_Get(t *testing.T) {

	const existingAccount = "f2d869a8-e5bc-4fbf-ad71-0000000000000"
	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-111111111111"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-2222222222222"

	const todoListIdWithoutElements = "f2d869a8-e5bc-4fbf-ad71-333333333333"
	const todoListIdWithElements = "f2d869a8-e5bc-4fbf-ad71-444444444444"
	const nonExistingTodoListId = "f2d869a8-e5bc-4fbf-ad71-555555555555"

	tests := []getListTestCase{
		{
			description:  "Invalid access token",
			accessToken:  nonExistingAccessToken,
			todoListId:   todoListIdWithoutElements,
			responseCode: http.StatusUnauthorized,
			responseBody: `{"error":"invalid access token"}`,
		},
		{
			description:  "Invalid todo list id parameter",
			accessToken:  validAccessToken,
			todoListId:   `invalid-list-id`,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"invalid todo list"}`,
		},
		{
			description:  "Todo list does not exist",
			accessToken:  validAccessToken,
			todoListId:   nonExistingTodoListId,
			responseCode: http.StatusBadRequest,
			responseBody: `{"error":"todo list does not exist"}`,
		},
		{
			description:  "Get empty todo list",
			accessToken:  validAccessToken,
			todoListId:   todoListIdWithoutElements,
			responseCode: http.StatusOK,
			responseBody: `{"todos":[]}`,
		},
		{
			description:  "Get todo list",
			accessToken:  validAccessToken,
			todoListId:   todoListIdWithElements,
			responseCode: http.StatusOK,
			responseBody: `{"todos":[{"created_by":"test user","description":"first todo","status":"todo","updated_at":"2024-01-01T00:00:00+00:00"},{"created_by":"test user","description":"second todo","status":"ongoing","updated_at":"2023-01-01T00:00:00+00:00"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			database := db.InMemoryDatabase()
			database.Users[existingAccount] = "test user"
			database.AccessTokens[validAccessToken] = existingAccount
			database.TodoLists[todoListIdWithoutElements] = []models.TodoDatabaseItem{}
			database.TodoLists[todoListIdWithElements] = []models.TodoDatabaseItem{
				{Description: "first todo", Status: "todo", User: existingAccount, UpdatedAt: todoTime(2024, 1, 1)},
				{Description: "second todo", Status: "ongoing", User: existingAccount, UpdatedAt: todoTime(2023, 1, 1)},
			}

			todoList := TodoLists{
				Database:     database,
				GenerateUuid: func() string { return "static_uuid" },
			}

			request := httptest.NewRequest(http.MethodGet, "/todolists", nil)
			request.SetPathValue("list_id", tt.todoListId)
			request.Header.Set("Authorization", tt.accessToken)
			writer := httptest.NewRecorder()

			todoList.Get(writer, request)

			assert.Equal(t, tt.responseCode, writer.Code)
			assert.Equal(t, tt.responseBody, writer.Body.String())
		})
	}
}

func todoTime(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.FixedZone("CEST", 1))
}
