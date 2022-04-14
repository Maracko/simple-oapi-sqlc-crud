//go:generate go run github.com/discord-gophers/goapi-gen --package=api --generate types,server,spec -o todo.gen.go ../todo-spec.yaml

package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/maracko/oapi-sqlc-crud/db"
)

type TodoServer struct {
	queries *db.Queries
}

func NewTodoServer(queries *db.Queries) *TodoServer {
	return &TodoServer{queries}
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*TodoServer)(nil)

func (s *TodoServer) GetTodos(w http.ResponseWriter, r *http.Request, params GetTodosParams) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Get by tags
	if params.Tags != nil {
		res, err := s.queries.ListTodosWithTags(ctx, params.Tags)
		if err != nil {
			render.Render(w, r, GetTodosJSONDefaultResponse(Error{Message: err.Error()}).Status(http.StatusInternalServerError))
			return
		}
		render.Render(w, r, GetTodosJSON200Response(convertDBTodosIntoTodos(res)))
		return
	}
	// Get all
	res, err := s.queries.ListTodos(ctx)
	if err != nil {
		render.Render(w, r, GetTodosJSONDefaultResponse(Error{Message: err.Error()}).Status(http.StatusInternalServerError))
		return
	}
	render.Render(w, r, GetTodosJSON200Response(convertDBTodosIntoTodos(res)))
}

func (s *TodoServer) AddTodo(w http.ResponseWriter, r *http.Request) {

	var newTodo AddTodoJSONRequestBody
	//Error handling probably isn't needed since the middleware takes care of it, but left for reference
	if err := render.Bind(r, &newTodo); err != nil {
		render.Render(
			w, r,
			AddTodoJSONDefaultResponse(Error{Message: "Invalid body format for todo"}).Status(http.StatusBadRequest),
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	todo, err := s.queries.CreateTodo(ctx, convertNewTodoIntoCreateTodoParams(newTodo))
	if err != nil {
		render.Render(
			w, r,
			AddTodoJSONDefaultResponse(Error{Message: err.Error()}).Status(http.StatusInternalServerError),
		)
		return
	}

	render.Render(
		w, r,
		AddTodoJSON200Response(convertDBTodoIntoTodo(todo)),
	)
}

func (s *TodoServer) FindTodoByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	todo, err := s.queries.GetTodo(ctx, id)
	if err != nil {
		render.Render(
			w, r,
			AddTodoJSONDefaultResponse(Error{Message: "cannot find todo"}).Status(http.StatusBadRequest),
		)
		return
	}

	render.Render(
		w, r,
		FindTodoByIDJSON200Response(convertDBTodoIntoTodo(todo)),
	)

}

func (s *TodoServer) DeleteTodo(w http.ResponseWriter, r *http.Request, id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := s.queries.GetTodo(ctx, id); err != nil {
		render.Render(
			w, r,
			DeleteTodoJSONDefaultResponse(Error{Message: "cannot find todo"}).Status(http.StatusBadRequest),
		)
		return
	}
	err := s.queries.DeleteTodo(ctx, id)
	if err != nil {
		render.Render(
			w, r,
			AddTodoJSONDefaultResponse(Error{Message: fmt.Sprint("cannot delete todo: ", err)}).Status(http.StatusBadRequest),
		)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
