package main

import (
	"fmt"
	"log"
	"net/http"

	"gql"
	"postgres"
	"server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
)

func main() {
	// Initialize our api and return a point to router for http
	// and a pointer to database
	router, db := initializeAPI()
	defer db.Close()

	// Listen on port 4000 and log all errors
	log.Fatal(http.ListenAndServe(":4000", router))
}

func initializeAPI() (*chi.Mux, *postgres.Db) {
	// Create a route
	router := chi.NewRouter()

	// Create a new connection to pg database
	db, err := postgres.New(
		postgres.ConnString("127.0.0.1", 5432, "postgres", "go_graphql_db"),
	)

	// If has a error log  will show
	if err != nil {
		log.Fatal(err)
	}

	// Create root query for graphql
	rootQuery := gql.NewRoot(db)

	// Create a graphql schema
	source, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// Create a server struct that holds a poiter to db
	s := server.Server{
		GqlSchema: &source,
	}

	// Add middleware to router
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
	)

	// Create graphql route
	router.Post("/graphql", s.GraphQL())

	return router, db
}
