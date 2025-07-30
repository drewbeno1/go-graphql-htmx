package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/drewbeno1/go-graphql-htmx/graph"
  "github.com/drewbeno1/go-graphql-htmx/graph/generated"
)

var resolver = &graph.Resolver{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: resolver},
		),
	)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	// Serve GraphQL endpoint to HTMX and Playground
	http.Handle("/query", srv)

	log.Printf("🚀 running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
