package storage

import (
	"testing"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupMaterialGORM(t *testing.T) *AlamcenSQlite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Material{})
	assert.NoError(t, err)
	return NewAlmacenSQLite(db)
}

func TestMaterialGORM_CrearYListar(t *testing.T) {
	repo := setupMaterialGORM(t)

	in := models.EntradaMaterial{
		Nombre:           "Cemento",
		Descripcion:      "Saco 50kg",
		Unidad:           "unidad",
		PrecioReferencia: "25.50",
	}
	mat, err := repo.CrearMateriales(in)
	assert.NoError(t, err)
	assert.Equal(t, "Cemento", mat.Nombre)
	assert.Equal(t, "unidad", mat.Unidad)
	assert.True(t, mat.PrecioReferencia.GreaterThan(decimal.Zero))
	assert.Greater(t, mat.ID, 0)

	todos := repo.ListarMateriales()
	assert.Len(t, todos, 1)
	assert.Equal(t, "Cemento", todos[0].Nombre)
	assert.Equal(t, mat.ID, todos[0].ID)
}

func TestMaterialGORM_ObtenerPorID(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})
	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Arena", Unidad: "m³", PrecioReferencia: "22.00"})

	mat, ok := repo.ObtenerMateriales(1)
	assert.True(t, ok)
	assert.Equal(t, "Cemento", mat.Nombre)

	_, ok = repo.ObtenerMateriales(99)
	assert.False(t, ok)
}

func TestMaterialGORM_Actualizar(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})

	act, ok := repo.ActualizarMateriales(1, models.EntradaMaterial{Nombre: "Cemento Plus", Descripcion: "Saco 50kg", Unidad: "unidad", PrecioReferencia: "30.00"})
	assert.True(t, ok)
	assert.Equal(t, "Cemento Plus", act.Nombre)

	_, ok = repo.ActualizarMateriales(99, models.EntradaMaterial{Nombre: "X", Unidad: "unidad", PrecioReferencia: "1"})
	assert.False(t, ok)
}

func TestMaterialGORM_Eliminar(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: "25.50"})

	assert.True(t, repo.EliminarMateriales(1))
	assert.False(t, repo.EliminarMateriales(1))
	assert.False(t, repo.EliminarMateriales(99))
}

// ─────────────────────────────────────────────
// helpers por entidad
// ─────────────────────────────────────────────

func setupManoObraGORM(t *testing.T) *AlamcenSQlite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{}})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&models.ManoObra{}))
	return NewAlmacenSQLite(db)
}

func setupEquipoGORM(t *testing.T) *AlamcenSQlite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{}})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&models.Equipo{}))
	return NewAlmacenSQLite(db)
}

func setupPrecioGORM(t *testing.T) *AlamcenSQlite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{}})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(
		&models.Material{},
		&models.ManoObra{},
		&models.Equipo{},
		&models.PrecioRecurso{},
	))
	return NewAlmacenSQLite(db)
}

func setupUsuarioGORM(t *testing.T) *UsuarioGORM {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{}})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&models.Usuario{}))
	return NewUsuarioRepository(db)
}

// =========================================================
// Mano de obra
// =========================================================

func TestManoObra_Listar(t *testing.T) {
	repo := setupManoObraGORM(t)
	_, _ = repo.CrearManoObra(models.EntradaManoObra{
		Descripcion:     "Maestro de obra",
		Categoria:       "oficial",
		Unidad:          "día",
		CostoReferencia: decimal.RequireFromString("35.00"),
	})
	assert.Len(t, repo.ListarManoObra(), 1)
}

func TestManoObra_CRUD(t *testing.T) {
	repo := setupManoObraGORM(t)

	creado, err := repo.CrearManoObra(models.EntradaManoObra{
		Descripcion:     "Maestro de obra",
		Categoria:       "oficial",
		Unidad:          "día",
		CostoReferencia: decimal.RequireFromString("35.00"),
	})
	assert.NoError(t, err)
	assert.Greater(t, creado.ID, 0)

	obt, ok := repo.ObtenerManoObra(creado.ID)
	assert.True(t, ok)
	assert.Equal(t, "oficial", obt.Categoria)

	act, ok := repo.ActualizarManoObra(creado.ID, models.EntradaManoObra{
		Descripcion:     "Maestro de obra senior",
		Categoria:       "oficial",
		Unidad:          "día",
		CostoReferencia: decimal.RequireFromString("40.00"),
	})
	assert.True(t, ok)
	assert.Equal(t, "Maestro de obra senior", act.Descripcion)

	assert.True(t, repo.EliminarManoObra(creado.ID))
	_, ok = repo.ObtenerManoObra(creado.ID)
	assert.False(t, ok)
}

// =========================================================
// Equipos
// =========================================================

func TestEquipo_CRUD(t *testing.T) {
	repo := setupEquipoGORM(t)

	creado, err := repo.CrearEquipo(models.EntradaEquipo{
		Nombre:     "Excavadora CAT 320",
		Tipo:       "pesado",
		Unidad:     "hora",
		CostoHora:  decimal.RequireFromString("85.00"),
		Disponible: true,
	})
	assert.NoError(t, err)
	assert.Greater(t, creado.ID, 0)
	assert.True(t, creado.Disponible)

	obt, err := repo.ObtenerEquipo(creado.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Excavadora CAT 320", obt.Nombre)

	act, err := repo.ActualizarEquipo(creado.ID, models.EntradaEquipo{
		Nombre:     "Excavadora CAT 320",
		Tipo:       "pesado",
		Unidad:     "hora",
		CostoHora:  decimal.RequireFromString("90.00"),
		Disponible: false,
	})
	assert.NoError(t, err)
	assert.False(t, act.Disponible)

	assert.NoError(t, repo.EliminarEquipo(creado.ID))
	assert.NoError(t, repo.EliminarEquipo(creado.ID))
}

// =========================================================
// Precios
// =========================================================

func TestPrecio_CRUD(t *testing.T) {
	repo := setupPrecioGORM(t)
	base := time.Now().UTC()

	// El recurso referenciado debe existir para que ExisteRecurso sea coherente.
	mat, err := repo.CrearMateriales(models.EntradaMaterial{
		Nombre:           "Cemento",
		Unidad:           "unidad",
		PrecioReferencia: "9.50",
	})
	assert.NoError(t, err)

	p, err := repo.CrearPrecio(models.EntradaPrecioRecurso{
		RecursoTipo:   "material",
		RecursoID:     mat.ID,
		Precio:        decimal.RequireFromString("9.50"),
		FechaVigencia: base,
	})
	assert.NoError(t, err)
	assert.Greater(t, p.ID, 0)

	obt, err := repo.ObtenerPrecio(p.ID)
	assert.NoError(t, err)
	assert.Equal(t, p.ID, obt.ID)

	assert.Len(t, repo.ListarPrecios(), 1)

	assert.Len(t, repo.HistorialPrecios("material", mat.ID), 1)

	vig, err := repo.PrecioVigente("material", mat.ID)
	assert.NoError(t, err)
	assert.Equal(t, "9.5", vig.Precio.String())

	assert.NoError(t, repo.ExisteRecurso("material", mat.ID))
	assert.Error(t, repo.ExisteRecurso("material", 99999))

	assert.NoError(t, repo.EliminarPrecio(p.ID))
	_, err = repo.ObtenerPrecio(p.ID)
	assert.Error(t, err)
}

// =========================================================
// Usuarios
// =========================================================

func TestUsuario_CRUD(t *testing.T) {
	repo := setupUsuarioGORM(t)

	creado, err := repo.CrearUsuario(models.EntradaUsuario{
		Email:    "user@example.com",
		Password: "hash-seguro",
	})
	assert.NoError(t, err)
	assert.Greater(t, creado.ID, 0)

	u, ok := repo.BuscarUsuarioPorEmail("user@example.com")
	assert.True(t, ok)
	assert.Equal(t, "user@example.com", u.Email)

	_, ok = repo.BuscarUsuarioPorEmail("noexiste@example.com")
	assert.False(t, ok)

	assert.Len(t, repo.ListarUsuarios(), 1)

	obt, ok := repo.ObtenerUsuarioPorID(creado.ID)
	assert.True(t, ok)
	assert.Equal(t, creado.ID, obt.ID)
}
