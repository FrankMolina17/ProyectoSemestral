package handlers

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

type ServerC struct {
	ManoObra      *services.ManoObraServise
	Material      *services.MaterialService
	Equipos       *services.EquipoService
	Precios       *services.PreciosService
	Autenticacion *services.AutenticacionService
}

func NewServerC(
	manoObra *services.ManoObraServise,
	material *services.MaterialService,
	equipos *services.EquipoService,
	precios *services.PreciosService,
	autenticacion *services.AutenticacionService,
) *ServerC {
	return &ServerC{
		ManoObra:      manoObra,
		Material:      material,
		Equipos:       equipos,
		Precios:       precios,
		Autenticacion: autenticacion,
	}
}
