package service

import (
	"github.com/fir1/port/internal/port/repository"
	"github.com/fir1/port/internal/port/service"
	"go.uber.org/fx"
)

var FxProvide = fx.Provide(
	service.NewPortService,
	repository.NewPostRepositoryMemoryDB,
)
