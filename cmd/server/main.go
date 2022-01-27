// This is an example of implementing the Pet Store from the OpenAPI documentation
// found at:
// https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v3.0/petstore.yaml

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/discord-gophers/goapi-gen/pkg/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/maracko/oapi-sqlc-crud/api"
	"github.com/maracko/oapi-sqlc-crud/db"

	_ "github.com/lib/pq"
)

func main() {
	// Flags to run the server, for more serious use should have config file or env. variables
	port := flag.Int("port", 8080, "Port for HTTP server")
	pgUser := flag.String("pg-user", "", "Postgres user")
	pgPass := flag.String("pg-pass", "", "Postgres password")
	pgHost := flag.String("pg-host", "localhost", "Postgres host")
	pgPort := flag.Int("pg-port", 5432, "Postgres port")
	pgDB := flag.String("pg-db", "todo-api", "Postgres database")
	help := flag.Bool("help", false, "Print help message")
	flag.Parse()

	if *pgUser == "" || *pgPass == "" || *help {
		printRequiredFlagsAndExit()
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	//Init the DB connection
	dbConn, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			*pgHost, *pgPort, *pgUser, *pgPass, *pgDB,
		),
	)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	err = dbConn.Ping()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	// Create an instance of queries
	queries := db.New(dbConn)
	// Create an instance of our handler which satisfies the generated interface
	todoServer := api.NewTodoServer(queries)

	// Create basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// Register todoServer above as the handler for the interface
	api.Handler(todoServer, api.WithRouter(r))

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}

func printRequiredFlagsAndExit() {
	msg := `To run the server you must specify the following flags:
	--pg-user Username for postgres server
	--pg-pass Password for postgres server
Following flags are optional:
	--pg-host Hostname for postgres server. Default localhost
	--pg-port Port for postgres server. Default 5432
	--pg-db Database name for postgres server. Default todo-api
	--port Port for HTTP server. Default 8080`

	fmt.Println(msg)
	os.Exit(1)
}
