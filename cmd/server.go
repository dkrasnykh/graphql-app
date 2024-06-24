package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"

	"github.com/dkrasnykh/graphql-app/graph"
	"github.com/dkrasnykh/graphql-app/internal/service"
	"github.com/dkrasnykh/graphql-app/internal/storage/database"
	"github.com/dkrasnykh/graphql-app/internal/storage/memory"
	"github.com/dkrasnykh/graphql-app/internal/subscription"
)

const (
	defaultPort = "8080"
	//databaseURL = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" // for local testing
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
	subscriptions := subscription.New()
	serv := service.New(storager, subscriptions)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Service:       serv,
		Subscriptions: subscriptions,
	}}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	//srv.AddTransport(&transport.Websocket{})

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
