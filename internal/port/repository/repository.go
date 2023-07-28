package repository

import (
	"context"

	"github.com/fir1/port/internal/port/model"
)

type PostRepositoryInterface interface {
	Create(ctx context.Context, key string, entity model.Port) error
	Update(ctx context.Context, key string, entity model.Port) error
	Get(ctx context.Context, key string) (model.Port, error)
	ListAll(ctx context.Context) (map[string]model.Port, error)
}
