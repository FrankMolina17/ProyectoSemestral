package repository_test

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) repository.ProformaRepository {
	t.Helper()
	db, err := repository.NuevaConexion("file::memory:?cache=shared")
	require.NoError(t, err)
	return repository.NuevoProformaRepository(db)
}

func TestProformaRepository_Actualizar(t *testing.T) {
	repo := setupRepo(t)

	creada, err := repo.CrearProforma(models.Proforma{Nombre: "Original", ObraID: 1})
	require.NoError(t, err)

	act, err := repo.ActualizarProforma(creada.ID, models.Proforma{Nombre: "Actualizado", PctGanancia: 0.15})
	require.NoError(t, err)
	assert.Equal(t, "Actualizado", act.Nombre)
	assert.Equal(t, 0.15, act.PctGanancia)

	_, err = repo.ActualizarProforma(999, models.Proforma{Nombre: "Nope"})
	assert.Error(t, err)
}

func TestProformaRepository_Eliminar(t *testing.T) {
	repo := setupRepo(t)

	creada, err := repo.CrearProforma(models.Proforma{Nombre: "Para borrar", ObraID: 1})
	require.NoError(t, err)

	assert.NoError(t, repo.EliminarProforma(creada.ID))
	assert.Error(t, repo.EliminarProforma(creada.ID))

	_, err = repo.ObtenerPorID(creada.ID)
	assert.Error(t, err)
}

func TestProformaRepository_Aprobar(t *testing.T) {
	repo := setupRepo(t)

	creada, err := repo.CrearProforma(models.Proforma{Nombre: "Proforma", ObraID: 1})
	require.NoError(t, err)

	aprobada, err := repo.AprobarProforma(creada.ID)
	require.NoError(t, err)
	assert.Equal(t, "aprobada", aprobada.Estado)

	_, err = repo.AprobarProforma(999)
	assert.Error(t, err)
}

func TestProformaRepository_Items(t *testing.T) {
	repo := setupRepo(t)

	proforma, err := repo.CrearProforma(models.Proforma{Nombre: "Con items", ObraID: 1})
	require.NoError(t, err)

	item, err := repo.AgregarItem(models.ProformaItem{
		ProformaID:     proforma.ID,
		TipoRecurso:    "material",
		Descripcion:    "Cemento",
		Cantidad:       10,
		PrecioPromedio: 12.50,
		Subtotal:       125.0,
	})
	require.NoError(t, err)
	assert.Greater(t, item.ID, 0)

	items, err := repo.ObtenerItems(proforma.ID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Cemento", items[0].Descripcion)
}

func TestProformaRepository_Notas(t *testing.T) {
	repo := setupRepo(t)

	proforma, err := repo.CrearProforma(models.Proforma{Nombre: "Con notas", ObraID: 1})
	require.NoError(t, err)

	nota, err := repo.AgregarNota(models.NotaProforma{
		ProformaID: proforma.ID,
		Contenido:  "Nota de prueba",
	})
	require.NoError(t, err)
	assert.Greater(t, nota.ID, 0)

	notas, err := repo.ObtenerNotas(proforma.ID)
	require.NoError(t, err)
	assert.Len(t, notas, 1)
	assert.Equal(t, "Nota de prueba", notas[0].Contenido)
}

func TestProformaRepository_Clientes(t *testing.T) {
	repo := setupRepo(t)

	cliente, err := repo.CrearCliente(models.Cliente{Nombre: "Juan", Ruc: "1234567890"})
	require.NoError(t, err)
	assert.Greater(t, cliente.ID, 0)

	clientes, err := repo.ObtenerClientes()
	require.NoError(t, err)
	assert.Len(t, clientes, 1)

	obt, err := repo.ObtenerClientePorID(cliente.ID)
	require.NoError(t, err)
	assert.Equal(t, "Juan", obt.Nombre)

	act, err := repo.ActualizarCliente(cliente.ID, models.Cliente{Nombre: "Juan Actualizado", Ruc: "1234567890"})
	require.NoError(t, err)
	assert.Equal(t, "Juan Actualizado", act.Nombre)

	assert.NoError(t, repo.EliminarCliente(cliente.ID))
	_, err = repo.ObtenerClientePorID(cliente.ID)
	assert.Error(t, err)
}
