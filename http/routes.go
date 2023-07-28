package http

func (s *Service) routes() {
	s.router.Get("/health", s.GetHealth)
	s.router.Get("/ports", s.listPorts)
	s.router.Post("/ports", s.savePorts)
	s.router.Post("/ports/from-file", s.savePortsFromFile)
}
