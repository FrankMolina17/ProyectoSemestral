package fakes

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNuevoProformaRepositoryFake(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	require.NotNil(t, f)
}

func TestProformaRepositoryFake_CrearYObtenerProforma(t *testing.T) {
	f := NuevoProformaRepositoryFake()

	p := models.Proforma{Nombre: "Proforma Test", ClienteID: 1}
	creada, err := f.CrearProforma(p)
	require.NoError(t, err)
	assert.Equal(t, 1, creada.ID)
	assert.Equal(t, "borrador", creada.Estado)

	obtenida, err := f.ObtenerPorID(creada.ID)
	require.NoError(t, err)
	assert.Equal(t, "Proforma Test", obtenida.Nombre)

	_, err = f.ObtenerPorID(999)
	assert.Error(t, err)
}

func TestProformaRepositoryFake_ObtenerTodos(t *testing.T) {
	f := NuevoProformaRepositoryFake()

	f.CrearProforma(models.Proforma{Nombre: "P1"})
	f.CrearProforma(models.Proforma{Nombre: "P2"})

	lista, err := f.ObtenerTodos()
	require.NoError(t, err)
	assert.Len(t, lista, 2)
}

func TestProformaRepositoryFake_ActualizarProforma(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	creada, _ := f.CrearProforma(models.Proforma{Nombre: "Original", PctGanancia: 0.1})

	actualizada, err := f.ActualizarProforma(creada.ID, models.Proforma{Nombre: "Actualizada", PctGanancia: 0.2})
	require.NoError(t, err)
	assert.Equal(t, "Actualizada", actualizada.Nombre)

	_, err = f.ActualizarProforma(999, models.Proforma{})
	assert.Error(t, err)
}

func TestProformaRepositoryFake_EliminarProforma(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	creada, _ := f.CrearProforma(models.Proforma{Nombre: "Para Eliminar"})

	err := f.EliminarProforma(creada.ID)
	assert.NoError(t, err)

	err = f.EliminarProforma(creada.ID)
	assert.Error(t, err)
}

func TestProformaRepositoryFake_AprobarProforma(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	creada, _ := f.CrearProforma(models.Proforma{Nombre: "Para Aprobar"})

	aprobada, err := f.AprobarProforma(creada.ID)
	require.NoError(t, err)
	assert.Equal(t, "aprobada", aprobada.Estado)

	_, err = f.AprobarProforma(999)
	assert.Error(t, err)
}

func TestProformaRepositoryFake_Items(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	creada, _ := f.CrearProforma(models.Proforma{Nombre: "Con Items"})

	item, err := f.AgregarItem(models.ProformaItem{ProformaID: creada.ID, Cantidad: 5, PrecioPromedio: 10.0})
	require.NoError(t, err)
	assert.Equal(t, 1, item.ID)
	assert.Equal(t, 50.0, item.Subtotal)

	_, err = f.AgregarItem(models.ProformaItem{ProformaID: 999})
	assert.Error(t, err)

	items, err := f.ObtenerItems(creada.ID)
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestProformaRepositoryFake_Notas(t *testing.T) {
	f := NuevoProformaRepositoryFake()
	creada, _ := f.CrearProforma(models.Proforma{Nombre: "Con Notas"})

	nota, err := f.AgregarNota(models.NotaProforma{ProformaID: creada.ID, Contenido: "Nota de prueba"})
	require.NoError(t, err)
	assert.Equal(t, 1, nota.ID)

	_, err = f.AgregarNota(models.NotaProforma{ProformaID: 999})
	assert.Error(t, err)

	notas, err := f.ObtenerNotas(creada.ID)
	require.NoError(t, err)
	assert.Len(t, notas, 1)
}

func TestProformaRepositoryFake_Clientes(t *testing.T) {
	f := NuevoProformaRepositoryFake()

	c, err := f.CrearCliente(models.Cliente{Nombre: "Cliente Test"})
	require.NoError(t, err)
	assert.Equal(t, 1, c.ID)

	clientes, err := f.ObtenerClientes()
	require.NoError(t, err)
	assert.Len(t, clientes, 1)

	obtenido, err := f.ObtenerClientePorID(c.ID)
	require.NoError(t, err)
	assert.Equal(t, "Cliente Test", obtenido.Nombre)

	_, err = f.ObtenerClientePorID(999)
	assert.Error(t, err)

	actualizado, err := f.ActualizarCliente(c.ID, models.Cliente{Nombre: "Actualizado"})
	require.NoError(t, err)
	assert.Equal(t, "Actualizado", actualizado.Nombre)

	_, err = f.ActualizarCliente(999, models.Cliente{})
	assert.Error(t, err)

	err = f.EliminarCliente(c.ID)
	assert.NoError(t, err)

	err = f.EliminarCliente(c.ID)
	assert.Error(t, err)
}
