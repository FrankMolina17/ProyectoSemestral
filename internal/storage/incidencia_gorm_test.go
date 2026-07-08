package storage

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIncidenciaGORM_CrearYBuscar(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.AutoMigrate(&models.Incidencia{})

	repo := NuevoAlmacenSQLite(db)

	incidencia := models.Incidencia{
		EntidadTipo:   "obra",
		EntidadID:     42,
		ResponsableID: 7,
		Titulo:        "Prueba de integración GORM",
		Descripcion:   "Test real contra SQLite :memory:",
		Estado:        "Pendiente",
	}

	creada := repo.CrearIncidencia(incidencia)
	require.NotZero(t, creada.ID, "Debe generar un ID")

	encontrada, ok := repo.BuscarIncidenciaPorID(creada.ID)
	require.True(t, ok)
	require.Equal(t, creada.ID, encontrada.ID)
	require.Equal(t, "Prueba de integración GORM", encontrada.Titulo)
}
