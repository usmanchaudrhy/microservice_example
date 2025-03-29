package memory

import (
	"awesomeProject/metadata/internal/repository"
	model "awesomeProject/metadata/pkg"
	"context"
	"sync"
)

// Repository defines a memory movie metadata repository
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

// Get function retrieves movie metadata for a movie by ID
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.Lock()
	defer r.Unlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

// Put adds movie metadata for a given movie id
func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = metadata
	return nil
}
