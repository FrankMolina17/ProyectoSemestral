package storage

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

type MaterialRepository interface {
	CrearMateriales(in models.EntradaMaterial) (*models.Material, error)
	ObtenerMateriales(id int) (*models.Material, bool)
	ListarMateriales() []*models.Material
	ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool)
	EliminarMateriales(id int) bool
}

type ManoObraRepository interface {
	CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error)
	ObtenerManoObra(id int) (*models.ManoObra, bool)
	ListarManoObra() []*models.ManoObra
	ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool)
	EliminarManoObra(id int) bool
}

type EquipoRepository interface {
	CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error)
	ObtenerEquipo(id int) (*models.Equipo, error)
	ListarEquipos() []*models.Equipo
	ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error)
	EliminarEquipo(id int) error
}

type PrecioRecursoRepository interface {
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
	ListarUsuarios() []*models.Usuario
	ObtenerUsuarioPorID(id int) (*models.Usuario, bool)
}

type IncidenciaRepository interface {
	ListarIncidencias() []models.Incidencia
	BuscarIncidenciaPorID(id int) (models.Incidencia, bool)
	BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool)
	CrearIncidencia(c models.Incidencia) models.Incidencia
	ActualizarIncidencia(id int, datos models.Incidencia) (models.Incidencia, bool)
	BorrarIncidencia(id int) bool
}

type ObraRepository interface {
	ListarObras() []models.Obra
	BuscarObraPorID(id int) (models.Obra, bool)
	CrearObra(o models.Obra) models.Obra
	ActualizarObra(id int, datos models.Obra) (models.Obra, bool)
	BorrarObra(id int) bool
}

type UserRepository interface {
	CrearUsuario(u models.Usuario) (models.Usuario, error)
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}

type Almacen interface {
	IncidenciaRepository
	ObraRepository
	UserRepository
}

// Chequeo en tiempo de compilación: si Memoria dejara de cumplir Almacen,
// el proyecto NO compila. Red de seguridad opcional.
var _ Almacen = (*Memoria)(nil)
