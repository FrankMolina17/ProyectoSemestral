package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetObrasGlobals() {
	storage.Obras = nil
	storage.ObraIDCounter = 1
}

func resetIncidenciasGlobals() {
	storage.Incidencias = nil
	storage.IncidenciaIDCounter = 1
}

func setupObrasRouter() chi.Router {
	resetObrasGlobals()
	r := chi.NewRouter()
	r.Post("/obras", handlers.CrearObraHandler)
	r.Get("/obras", handlers.ObtenerObrasHandler)
	r.Get("/obras/{id}", handlers.ObtenerObraHandler)
	r.Put("/obras/{id}", handlers.ActualizarObraHandler)
	r.Delete("/obras/{id}", handlers.EliminarObraHandler)
	return r
}

func setupIncidenciasRouter() chi.Router {
	resetIncidenciasGlobals()
	r := chi.NewRouter()
	r.Post("/incidencias", handlers.CrearIncidenciaHandler)
	r.Get("/incidencias", handlers.ObtenerIncidenciasHandler)
	r.Get("/incidencias/{id}", handlers.ObtenerIncidenciaPorIDHandler)
	r.Get("/incidencias/por/{tipo}/{id}", handlers.ObtenerIncidenciasPorEntidadHandler)
	r.Put("/incidencias/{id}", handlers.ActualizarIncidenciaHandler)
	r.Delete("/incidencias/{id}", handlers.EliminarIncidenciaHandler)
	return r
}

// ─────────────────────────────────────────────
// OBRAS
// ─────────────────────────────────────────────

func TestObraHandler_CrearYObtener(t *testing.T) {
	r := setupObrasRouter()

	t.Run("crear obra -> 201", func(t *testing.T) {
		body := `{"nombre":"Edificio Central","user_id":1}`
		req := httptest.NewRequest(http.MethodPost, "/obras", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var obra map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &obra))
		assert.Equal(t, "Edificio Central", obra["nombre"])
		assert.Equal(t, "planificacion", obra["estado"])
	})

	t.Run("listar obras -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/obras", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var obras []interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &obras))
		assert.Len(t, obras, 1)
	})

	t.Run("obtener por id -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/obras/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener id inexistente -> 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/obras/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("id invalido -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/obras/abc", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestObraHandler_ActualizarYEliminar(t *testing.T) {
	r := setupObrasRouter()

	body := `{"nombre":"Obra Inicial","user_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/obras", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("actualizar -> 200", func(t *testing.T) {
		body := `{"nombre":"Obra Actualizada"}`
		req := httptest.NewRequest(http.MethodPut, "/obras/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("actualizar id inexistente -> 400", func(t *testing.T) {
		body := `{"nombre":"Nope"}`
		req := httptest.NewRequest(http.MethodPut, "/obras/999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("actualizar con id invalido -> 400", func(t *testing.T) {
		body := `{"nombre":"Nope"}`
		req := httptest.NewRequest(http.MethodPut, "/obras/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/obras/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("eliminar inexistente -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/obras/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("eliminar con id invalido -> 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/obras/abc", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestObraHandler_CrearConNombreVacio(t *testing.T) {
	r := setupObrasRouter()

	body := `{"nombre":"","user_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/obras", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ─────────────────────────────────────────────
// INCIDENCIAS
// ─────────────────────────────────────────────

func TestIncidenciaHandler_CrearYObtener(t *testing.T) {
	r := setupIncidenciasRouter()

	t.Run("crear incidencia -> 201", func(t *testing.T) {
		body := `{"titulo":"Filtración","entidad_tipo":"obra","entidad_id":1,"prioridad":"alta"}`
		req := httptest.NewRequest(http.MethodPost, "/incidencias", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var inc map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &inc))
		assert.Equal(t, "Filtración", inc["titulo"])
		assert.Equal(t, "abierta", inc["estado"])
	})

	t.Run("crear con datos invalidos -> 400", func(t *testing.T) {
		body := `{invalid}`
		req := httptest.NewRequest(http.MethodPost, "/incidencias", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("listar incidencias -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener por id -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener id inexistente -> 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("obtener por entidad -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias/por/obra/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener por entidad inexistente -> 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias/por/obra/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestIncidenciaHandler_ActualizarYEliminar(t *testing.T) {
	r := setupIncidenciasRouter()

	body := `{"titulo":"Incidencia inicial","entidad_tipo":"obra","entidad_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/incidencias", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("actualizar -> 200", func(t *testing.T) {
		body := `{"titulo":"Actualizado","estado":"cerrada"}`
		req := httptest.NewRequest(http.MethodPut, "/incidencias/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("actualizar inexistente -> 404", func(t *testing.T) {
		body := `{"titulo":"Nope"}`
		req := httptest.NewRequest(http.MethodPut, "/incidencias/999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/incidencias/1", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("eliminar inexistente -> 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/incidencias/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestIncidenciaHandler_Vacia(t *testing.T) {
	r := setupIncidenciasRouter()

	t.Run("listar vacio -> 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/incidencias", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
