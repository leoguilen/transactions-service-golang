package inmemory

import (
	"sync"

	"github.com/leoguilen/transactions/internal/domain"
)

// TodoRepository is an in-memory repository for todos used in tests
type TodoRepository struct {
	mu    sync.RWMutex
	todos map[string]*domain.Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{todos: make(map[string]*domain.Todo)}
}

func (r *TodoRepository) Save(todo *domain.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.todos[todo.ID] = todo
	return nil
}

func (r *TodoRepository) GetAll() ([]*domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Todo, 0, len(r.todos))
	for _, t := range r.todos {
		res = append(res, t)
	}
	return res, nil
}
