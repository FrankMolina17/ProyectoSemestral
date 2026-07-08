package storage

import "Sistem-Inte-Gestion-Control-Obras/internal/models"

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
