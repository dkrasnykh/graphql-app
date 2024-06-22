package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/dkrasnykh/graphql-app/graph"
	"github.com/dkrasnykh/graphql-app/internal/service"
	"github.com/dkrasnykh/graphql-app/internal/storage/database"
	"github.com/dkrasnykh/graphql-app/internal/storage/memory"
)

const (
	defaultPort = "8080"
	//databaseURL = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	databaseURL = "postgres://postgres:password@db:5432/postgres?sslmode=disable" // build docker image
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var isMemory bool
	flag.BoolVar(&isMemory, "m", false, "use memory storage")
	flag.Parse()

	storager := storage(isMemory)

	serv := service.New(storager)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Service: serv}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func storage(isMemory bool) service.Storager {
	if isMemory {
		return memory.New()
	}
	err := database.Migrate(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := database.New(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	return storage
}
