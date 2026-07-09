package services

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/stretchr/testify/assert"
)

func resetIncidencias() {
	storage.Incidencias = nil
	storage.IncidenciaIDCounter = 1
}

func TestCrearIncidencia(t *testing.T) {
	resetIncidencias()

	t.Run("valida -> success", func(t *testing.T) {
		inc, err := CrearIncidencia(models.Incidencia{
			Titulo:      "Filtración",
			EntidadTipo: "obra",
			EntidadID:   1,
			Prioridad:   "alta",
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, inc.ID)
		assert.Equal(t, "abierta", inc.Estado)
	})

	t.Run("titulo vacio -> error", func(t *testing.T) {
		_, err := CrearIncidencia(models.Incidencia{Titulo: "", EntidadTipo: "obra", EntidadID: 1})
		assert.Error(t, err)
	})

	t.Run("entidad invalida -> error", func(t *testing.T) {
		_, err := CrearIncidencia(models.Incidencia{Titulo: "Test", EntidadTipo: "", EntidadID: 0})
		assert.Error(t, err)
	})
}

func TestObtenerIncidencias(t *testing.T) {
	resetIncidencias()

	CrearIncidencia(models.Incidencia{Titulo: "Inc1", EntidadTipo: "obra", EntidadID: 1})
	CrearIncidencia(models.Incidencia{Titulo: "Inc2", EntidadTipo: "obra", EntidadID: 1})
	CrearIncidencia(models.Incidencia{Titulo: "Inc3", EntidadTipo: "proforma", EntidadID: 2})

	t.Run("todas", func(t *testing.T) {
		assert.Len(t, ObtenerIncidencias(), 3)
	})

	t.Run("por entidad", func(t *testing.T) {
		result := ObtenerIncidenciasPorEntidad("obra", 1)
		assert.Len(t, result, 2)
	})

	t.Run("por ID existente", func(t *testing.T) {
		inc, ok := ObtenerIncidenciaPorID(1)
		assert.True(t, ok)
		assert.Equal(t, "Inc1", inc.Titulo)
	})

	t.Run("por ID inexistente", func(t *testing.T) {
		_, ok := ObtenerIncidenciaPorID(999)
		assert.False(t, ok)
	})
}

func TestActualizarIncidencia(t *testing.T) {
	resetIncidencias()

	CrearIncidencia(models.Incidencia{Titulo: "Original", EntidadTipo: "obra", EntidadID: 1, Estado: "abierta"})

	t.Run("actualizar existente", func(t *testing.T) {
		inc, err := ActualizarIncidencia(1, models.Incidencia{Titulo: "Actualizado", Estado: "cerrada"})
		assert.NoError(t, err)
		assert.Equal(t, "Actualizado", inc.Titulo)
		assert.Equal(t, "cerrada", inc.Estado)
	})

	t.Run("actualizar inexistente -> error", func(t *testing.T) {
		_, err := ActualizarIncidencia(999, models.Incidencia{Titulo: "Nope"})
		assert.Error(t, err)
	})
}

func TestEliminarIncidencia(t *testing.T) {
	resetIncidencias()

	CrearIncidencia(models.Incidencia{Titulo: "Para borrar", EntidadTipo: "obra", EntidadID: 1})

	assert.NoError(t, EliminarIncidencia(1))
	assert.Error(t, EliminarIncidencia(1))
	assert.Len(t, ObtenerIncidencias(), 0)
}

func resetObras() {
	storage.Obras = nil
	storage.ObraIDCounter = 1
}

func TestCrearObra(t *testing.T) {
	resetObras()

	t.Run("valida -> success", func(t *testing.T) {
		obra, err := CrearObraServicio(models.Obra{Nombre: "Edificio", UserID: 1})
		assert.NoError(t, err)
		assert.Equal(t, 1, obra.ID)
		assert.Equal(t, "planificacion", obra.Estado)
	})

	t.Run("nombre vacio -> error", func(t *testing.T) {
		_, err := CrearObraServicio(models.Obra{Nombre: "", UserID: 1})
		assert.Error(t, err)
	})

	t.Run("user_id invalido -> error", func(t *testing.T) {
		_, err := CrearObraServicio(models.Obra{Nombre: "Obra", UserID: 0})
		assert.Error(t, err)
	})
}

func TestObtenerObras(t *testing.T) {
	resetObras()

	CrearObraServicio(models.Obra{Nombre: "Obra 1", UserID: 1})
	CrearObraServicio(models.Obra{Nombre: "Obra 2", UserID: 2})

	assert.Len(t, ObtenerObras(), 2)

	obra, ok := ObtenerObra(1)
	assert.True(t, ok)
	assert.Equal(t, "Obra 1", obra.Nombre)

	_, ok = ObtenerObra(999)
	assert.False(t, ok)
}

func TestActualizarYEliminarObra(t *testing.T) {
	resetObras()

	CrearObraServicio(models.Obra{Nombre: "Original", UserID: 1})

	act, ok := ActualizarObra(1, models.Obra{Nombre: "Actualizado"})
	assert.True(t, ok)
	assert.Equal(t, "Actualizado", act.Nombre)

	_, ok = ActualizarObra(999, models.Obra{Nombre: "Nope"})
	assert.False(t, ok)

	assert.True(t, EliminarObra(1))
	assert.False(t, EliminarObra(1))
	assert.Len(t, ObtenerObras(), 0)
}
