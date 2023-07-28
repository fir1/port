package repository

import (
	"context"
	"sync"

	"github.com/fir1/port/internal/port/model"
)

type PostRepositoryMemoryDB struct {
	// Initialize your database connection here
	storage map[string]model.Port
	mu      sync.RWMutex
}

func NewPostRepositoryMemoryDB() PostRepositoryInterface {
	return &PostRepositoryMemoryDB{
		storage: make(map[string]model.Port),
		mu:      sync.RWMutex{},
	}
}

func (r *PostRepositoryMemoryDB) Get(ctx context.Context, key string) (model.Port, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	port, found := r.storage[key]
	if !found {
		return model.Port{}, ErrObjectNotFound{}
	}
	return port, nil
}

func (r *PostRepositoryMemoryDB) Create(ctx context.Context, key string, entity model.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[key] = entity
	return nil
}

func (r *PostRepositoryMemoryDB) Update(ctx context.Context, key string, entity model.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[key] = entity
	return nil
}

func (r *PostRepositoryMemoryDB) ListAll(ctx context.Context) (map[string]model.Port, error) {
	return r.storage, nil
}
