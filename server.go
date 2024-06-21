package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/dkrasnykh/graphql-app/graph"
	"github.com/dkrasnykh/graphql-app/internal/service"
	"github.com/dkrasnykh/graphql-app/internal/storage/database"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	databaseURL := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	timeout := time.Second * 2

	err := database.Migrate(databaseURL, time.Second*2)
	if err != nil {
		panic(err)
	}

	postStorage, err := database.NewPostPostgres(databaseURL, timeout)
	if err != nil {
		panic(err)
	}

	commentStorage, err := database.NewCommentPostgres(databaseURL, timeout)
	if err != nil {
		panic(err)
	}

	servicePost := service.NewPostService(postStorage)
	serviceComment := service.NewComment(commentStorage)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Post: servicePost, Comment: serviceComment}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
