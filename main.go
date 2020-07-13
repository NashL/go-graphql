package main

import (
	"github.com/gorilla/websocket"
	dbConnection "github.com/nashl/online-store-server/database"
	"github.com/nashl/online-store-server/graph"
	"github.com/nashl/online-store-server/graph/generated"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/joho/godotenv/autoload"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

const defaultPort = "9990"

func main() {
	router := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	//start connecting to the database
	dbConnection.NewDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//srv := handler.New(starwars.NewExecutableSchema(starwars.NewResolver()))
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}

	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}



