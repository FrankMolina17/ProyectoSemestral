package handlers

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	
)

type ServerC struct {
	ManoObra *services.ManoObraServise
	Material *services.MaterialService
	Equipos *services.EquipoService
	Precios *services.PreciosService
	Autenticacion *services.AutenticacionService
}
func NewServerC(manoObra *services.ManoObraServise, material *services.MaterialService, equipos *services.EquipoService, precios *services.PreciosService, autenticacion *services.AutenticacionService) *ServerC {
	return &ServerC{
		ManoObra: manoObra,
		Material: material,
		Equipos: equipos,
		Precios: precios,
		Autenticacion: autenticacion,
	}

}


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

type Deps struct{
	ManoObra *services.ManoObraServise
	Material *services.MaterialService
	Equipos *services.EquipoService
	Precios *services.PreciosService
	Autenticacion *services.AutenticacionService
}

func NewServer (d Deps) *ServerC{
	return &ServerC{
		ManoObra: d.ManoObra,
		Material: d.Material,
		Equipos: d.Equipos,
		Precios: d.Precios,
		Autenticacion: d.Autenticacion,
	}
}
