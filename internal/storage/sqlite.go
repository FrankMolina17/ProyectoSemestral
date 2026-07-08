package storage

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"gorm.io/gorm"
)

// AlmacenSQLite implementa la interfaz Almacen usando GORM sobre SQLite.
//
// Fíjense: los métodos tienen EXACTAMENTE las mismas firmas que los de Memoria.
// Por eso el Server y los handlers no se enteran de cuál de los dos reciben.
type AlmacenSQLite struct {
	db *gorm.DB
}

// NuevoAlmacenSQLite envuelve una conexión *gorm.DB ya abierta.
func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

// =========================================================
// INCIDENCIAS
// =========================================================

func (a *AlmacenSQLite) ListarIncidencias() []models.Incidencia {
	var incidencia []models.Incidencia
	a.db.Find(&incidencia)
	return incidencia
}

func (a *AlmacenSQLite) BuscarIncidenciaPorID(id int) (models.Incidencia, bool) {
	var i models.Incidencia
	if err := a.db.First(&i, id).Error; err != nil {
		// Absorbemos el error de la DB y conservamos la firma comma-ok.
		return models.Incidencia{}, false
	}
	return i, true
}

func (a *AlmacenSQLite) BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool) {
	var i models.Incidencia
	if err := a.db.Where("entidad_tipo = ?", tipo).First(&i, id).Error; err != nil {
		// Absorbemos el error de la DB y conservamos la firma comma-ok.
		return models.Incidencia{}, false
	}
	return i, true
}

func (a *AlmacenSQLite) CrearIncidencia(i models.Incidencia) models.Incidencia {
	a.db.Create(&i) // GORM rellena el ID autogenerado en &p
	return i
}

func (a *AlmacenSQLite) ActualizarIncidencia(id int, datos models.Incidencia) (models.Incidencia, bool) {
	var existente models.Incidencia
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Incidencia{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarIncidencia(id int) bool {
	res := a.db.Delete(&models.Incidencia{}, id)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) ListarObras() []models.Obra {
	var obras []models.Obra
	a.db.Find(&obras)
	return obras
}

func (a *AlmacenSQLite) BuscarObraPorID(id int) (models.Obra, bool) {
	var o models.Obra
	if err := a.db.First(&o, id).Error; err != nil {
		return models.Obra{}, false
	}
	return o, true
}

func (a *AlmacenSQLite) CrearObra(o models.Obra) models.Obra {
	a.db.Create(&o)
	return o
}

func (a *AlmacenSQLite) ActualizarObra(id int, datos models.Obra) (models.Obra, bool) {
	var existente models.Obra
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Obra{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarObra(id int) bool {
	res := a.db.Delete(&models.Obra{}, id)
	return res.RowsAffected > 0
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario

	if err := a.db.First(&u, email).Error; err != nil {
		return models.Usuario{}, false
	}

	return u, true
}

func (a *AlmacenSQLite) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	a.db.Create(&u)
	return u, nil
}

// Chequeo en tiempo de compilación: AlmacenSQLite debe cumplir Almacen.
var _ Almacen = (*AlmacenSQLite)(nil)
