package service

import (
	"context"

	"github.com/fir1/port/internal/port/model"
)

func (s PortService) ListPorts(ctx context.Context) (map[string]model.Port, error) {
	return s.repository.ListAll(ctx)
}
