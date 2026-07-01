package handlers

import "Sistem-Inte-Gestion-Control-Obras/internal/services"

type Server struct {
	IncidenciaService *services.IncidenciaService
	Auth              *services.AuthService
}

func NewServer(incidencias *services.IncidenciaService,
	auth *services.AuthService) *Server {
	return &Server{
		IncidenciaService: incidencias,
		Auth:              auth,
	}
}
