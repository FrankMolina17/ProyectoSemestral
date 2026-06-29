package repository_test

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProformaRepository_CrearYObtenerTodos usa GORM con sqlite :memory:
// y verifica que crear una proforma se refleje al listar.
func TestProformaRepository_CrearYObtenerTodos(t *testing.T) {
	db, err := repository.NuevaConexion("file::memory:?cache=shared")
	require.NoError(t, err)

	repo := repository.NuevoProformaRepository(db)

	creada, err := repo.CrearProforma(models.Proforma{
		Nombre: "Proforma puente vehicular",
		ObraID: 42,
	})
	require.NoError(t, err)
	require.NotZero(t, creada.ID)
	assert.Equal(t, "borrador", creada.Estado)

	lista, err := repo.ObtenerTodos()
	require.NoError(t, err)
	require.Len(t, lista, 1)
	assert.Equal(t, creada.ID, lista[0].ID)
	assert.Equal(t, "Proforma puente vehicular", lista[0].Nombre)
	assert.Equal(t, 42, lista[0].ObraID)

	encontrada, err := repo.ObtenerPorID(creada.ID)
	require.NoError(t, err)
	assert.Equal(t, creada.Nombre, encontrada.Nombre)
}
