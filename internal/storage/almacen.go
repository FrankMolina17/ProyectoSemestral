package storage

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	
)


type MaterialRepository interface{
	CrearMateriales(in models.EntradaMaterial) (*models.Material, error)
	ObtenerMateriales(id int) (*models.Material, bool)
	ListarMateriales() []*models.Material
	ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool)
	EliminarMateriales(id int) bool
}

type ManoObraRepository interface{
	CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error)
	ObtenerManoObra(id int) (*models.ManoObra, bool)
	ListarManoObra() []*models.ManoObra
	ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool)
	EliminarManoObra(id int) bool
}

type EquipoRepository interface{
	CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error)
	ObtenerEquipo(id int) (*models.Equipo, error)
	ListarEquipos() []*models.Equipo
	ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error)
	EliminarEquipo(id int) error
}

type PrecioRecursoRepository interface{
	ListarPrecios() []*models.PrecioRecurso 
	ObtenerPrecio(id int) (*models.PrecioRecurso, error)
	CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error)
	HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso
	PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error)
	ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) 
	ExisteRecurso(tipo string, id int) error 
	EliminarPrecio(id int) error 
	
}

type UsuarioRepository interface {
	CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}