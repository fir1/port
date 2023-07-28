package service

import (
	"github.com/fir1/port/config"
	"github.com/fir1/port/internal/port/repository"
)

type PortService struct {
	repository repository.PostRepositoryInterface
	config     config.Config
}

func NewPortService(rp repository.PostRepositoryInterface,
	cnf config.Config) PortService {
	return PortService{
		repository: rp,
		config:     cnf,
	}
}
