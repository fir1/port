package http

import (
	"github.com/fir1/port/config"
	"github.com/fir1/port/internal/port/service"
	"github.com/fir1/port/pkg/cache"
	"github.com/go-chi/chi/v5"

	"github.com/sirupsen/logrus"
)

type Service struct {
	router            *chi.Mux
	logger            *logrus.Logger
	stockSymbol       string
	stockNumberOfDays int
	config            config.Config
	cacheClient       cache.CacheClientInterface
	portService       service.PortService
}

func NewService(logger *logrus.Logger,
	cnf config.Config,
	cc cache.CacheClientInterface,
	ps service.PortService,
) *Service {
	return &Service{
		logger:      logger,
		config:      cnf,
		cacheClient: cc,
		portService: ps,
	}
}
