package api

import (
	"database/sql"

	"github.com/maracko/oapi-sqlc-crud/db"
)

func convertDBTodosIntoTodos(todos []*db.Todo) []Todo {
	var res []Todo
	for _, t := range todos {
		var (
			content, title string
		)

		if t.Content.Valid {
			content = t.Content.String
		}
		if t.Title.Valid {
			title = t.Title.String
		}

		res = append(res, Todo{
			NewTodo: NewTodo{
				Content: content,
				Title:   &title,
				Tags:    t.Tags,
			},
			ID: t.ID,
		})
	}
	return res
}

func convertDBTodoIntoTodo(todo *db.Todo) Todo {
	var (
		t Todo
	)

	if todo.Content.Valid {
		t.Content = todo.Content.String
	}
	if todo.Title.Valid {
		t.Title = &todo.Title.String
	}

	t.Tags = todo.Tags
	t.ID = todo.ID

	return t
}

func convertNewTodoIntoCreateTodoParams(todo AddTodoJSONRequestBody) db.CreateTodoParams {
	var (
		content, title sql.NullString
	)

	if todo.Content != "" {
		content.String = todo.Content
		content.Valid = true
	}

	if todo.Title != nil {
		title.String = *todo.Title
		title.Valid = true
	}

	return db.CreateTodoParams{
		Title:   title,
		Tags:    todo.Tags,
		Content: content,
	}
}
