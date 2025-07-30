package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	tasks []*Task
	mu    sync.Mutex
)

// schema
var taskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Task",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"title":     &graphql.Field{Type: graphql.String},
		"completed": &graphql.Field{Type: graphql.Boolean},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"tasks": &graphql.Field{
			Type: graphql.NewList(taskType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				return tasks, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addTask": &graphql.Field{
			Type: taskType,
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				title := p.Args["title"].(string)
				task := &Task{ID: uuid.New().String(), Title: title}
				tasks = append(tasks, task)
				return task, nil
			},
		},
		"toggleTask": &graphql.Field{
			Type: taskType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mu.Lock()
				defer mu.Unlock()
				id := p.Args["id"].(string)
				for _, t := range tasks {
					if t.ID == id {
						t.Completed = !t.Completed
						return t, nil
					}
				}
				return nil, fmt.Errorf("not found")
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			Query string `json:"query"`
		}
		json.NewDecoder(r.Body).Decode(&params)
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: params.Query,
		})
		json.NewEncoder(w).Encode(result)
	})

	http.ListenAndServe(":8080", nil)
}
