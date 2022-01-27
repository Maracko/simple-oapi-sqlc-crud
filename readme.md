# Example API using sqlc and goapi-gen

API Created for testing purposes from an openapi specification using `goapi-gen`.  
CRUD methods generated using `sqlc`.  
It supports CRUD operations for "todo" entities.

## Modifying API schema

1. Modify the `todo-spec.yaml`.
2. Run `go generate ./...` inside root to regenerate the API.

**Important** If go generate doesn't work(`go mod tidy` removes it's dependancies) run `go get github.com/discord-gophers/goapi-gen` first.

## SQL

To generate SQL boilerplate for go you need `sqlc` installed.  
Install it using: `go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

## Modifying SQL schema

1. Add new migration files inside the `db/sql/schema` folder.
2. Add the path to new migration file to `db/sqlc.yaml` inside the schema array.
3. Run `sqlc generate` inside `db/sql` folder.

## Modifying SQL queries

1. Add or modify files inside the `db/sql/queries..` folder.
2. If you added a new queries file (e.g. a new entity), also add it to `db/sqlc.yaml` queries array.
3. Run `sqlc generate` inside `db/sql` folder.

## Viewing SWAGGER docs

Upload the contents of `todo-spec.yaml` to [swagger](https://editor.swagger.io/) to view the interactive API schema.

## Running the server

run/build the server inside `cmd/server` folder
