package main

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"

	"github.com/dkrasnykh/graphql-app/graph"
	"github.com/dkrasnykh/graphql-app/internal/config"
	"github.com/dkrasnykh/graphql-app/internal/service"
	"github.com/dkrasnykh/graphql-app/internal/storage/database"
	"github.com/dkrasnykh/graphql-app/internal/storage/memory"
	"github.com/dkrasnykh/graphql-app/internal/subscription"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	storager := storage(cfg.IsMemoty, cfg.DatabaseURL)
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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

func storage(isMemory bool, databaseURL string) service.Storager {
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
