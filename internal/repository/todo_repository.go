package repository

import "github.com/leoguilen/transactions/internal/domain"

// TodoRepository defines persistence operations for todos (port)
type TodoRepository interface {
	Save(todo *domain.Todo) error
	GetAll() ([]*domain.Todo, error)
}
