package graph

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/drewbeno1/go-graphql-htmx/graph/model"
)

type Resolver struct {
	mu    sync.Mutex
	Tasks []*model.Task
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Resolver.Tasks, nil // ✅ this works
}

func (r *mutationResolver) AddTask(ctx context.Context, title string) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task := &model.Task{
		ID:        uuid.New().String(),
		Title:     title,
		Completed: false,
	}
	r.Tasks = append(r.Tasks, task)
	return task, nil
}

func (r *mutationResolver) ToggleTask(ctx context.Context, id string) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range r.Tasks {
		if task.ID == id {
			task.Completed = !task.Completed
			return task, nil
		}
	}
	return nil, nil
}
