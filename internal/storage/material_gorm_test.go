package storage

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"github.com/glebarez/sqlite"
)

func setupMaterialGORM(t *testing.T) *MaterialGORM {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Material{})
	assert.NoError(t, err)
	return NewMaterialGORM(db)
}

func TestMaterialGORM_CrearYListar(t *testing.T) {
	repo := setupMaterialGORM(t)

	in := models.EntradaMaterial{
		Nombre:           "Cemento",
		Descripcion:      "Saco 50kg",
		Unidad:           "unidad",
		PrecioReferencia: decimal.NewFromFloat(25.50),
	}
	mat, err := repo.CrearMateriales(in)
	assert.NoError(t, err)
	assert.Equal(t, "Cemento", mat.Nombre)
	assert.Equal(t, "unidad", mat.Unidad)
	assert.Equal(t, decimal.NewFromFloat(25.50), mat.PrecioReferencia)
	assert.Greater(t, mat.ID, 0)

	todos := repo.ListarMateriales()
	assert.Len(t, todos, 1)
	assert.Equal(t, "Cemento", todos[0].Nombre)
	assert.Equal(t, mat.ID, todos[0].ID)
}

func TestMaterialGORM_ObtenerPorID(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: decimal.NewFromFloat(25.50)})
	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Arena", Unidad: "m³", PrecioReferencia: decimal.NewFromFloat(22.00)})

	mat, ok := repo.ObtenerMateriales(1)
	assert.True(t, ok)
	assert.Equal(t, "Cemento", mat.Nombre)

	_, ok = repo.ObtenerMateriales(99)
	assert.False(t, ok)
}

func TestMaterialGORM_Actualizar(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: decimal.NewFromFloat(25.50)})

	act, ok := repo.ActualizarMateriales(1, models.EntradaMaterial{Nombre: "Cemento Plus", Descripcion: "Saco 50kg", Unidad: "unidad", PrecioReferencia: decimal.NewFromFloat(30.00)})
	assert.True(t, ok)
	assert.Equal(t, "Cemento Plus", act.Nombre)
	assert.Equal(t, decimal.NewFromFloat(30.00), act.PrecioReferencia)

	_, ok = repo.ActualizarMateriales(99, models.EntradaMaterial{Nombre: "X", Unidad: "unidad", PrecioReferencia: decimal.NewFromFloat(1)})
	assert.False(t, ok)
}

func TestMaterialGORM_Eliminar(t *testing.T) {
	repo := setupMaterialGORM(t)

	_, _ = repo.CrearMateriales(models.EntradaMaterial{Nombre: "Cemento", Unidad: "unidad", PrecioReferencia: decimal.NewFromFloat(25.50)})

	assert.True(t, repo.EliminarMateriales(1))
	assert.False(t, repo.EliminarMateriales(1))
	assert.False(t, repo.EliminarMateriales(99))
}
