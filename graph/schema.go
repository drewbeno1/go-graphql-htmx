package graph

import (
	"sync"

	"go-graphql-htmx/models"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	mu  sync.Mutex
)

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Task{})
}

var TaskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Task",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"title":     &graphql.Field{Type: graphql.String},
		"completed": &graphql.Field{Type: graphql.Boolean},
	},
})

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"tasks": &graphql.Field{
			Type: graphql.NewList(TaskType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var tasks []models.Task
				if err := db.Find(&tasks).Error; err != nil {
					return nil, err
				}
				return tasks, nil
			},
		},
	},
})

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"addTask": &graphql.Field{
			Type: TaskType,
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				title := p.Args["title"].(string)
				task := models.Task{ID: uuid.New().String(), Title: title}
				if err := db.Create(&task).Error; err != nil {
					return nil, err
				}
				return task, nil
			},
		},
		"toggleTask": &graphql.Field{
			Type: TaskType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				var task models.Task
				if err := db.First(&task, "id = ?", id).Error; err != nil {
					return nil, err
				}
				task.Completed = !task.Completed
				if err := db.Save(&task).Error; err != nil {
					return nil, err
				}
				return task, nil
			},
		},
	},
})

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    RootQuery,
	Mutation: RootMutation,
})
