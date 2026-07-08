package handlers

import "Sistem-Inte-Gestion-Control-Obras/internal/services"

type Server struct {
	IncidenciaService *services.IncidenciaService
	ObraService       *services.ObraService
	Auth              *services.AuthService
}

type Deps struct {
	IncidenciaService *services.IncidenciaService
	ObraService       *services.ObraService
	Auth              *services.AuthService
}

func NewServer(d Deps) *Server {
	return &Server{
		IncidenciaService: d.IncidenciaService,
		Auth:              d.Auth,
	}
}
