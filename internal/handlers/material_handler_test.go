package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMaterialRouter(t *testing.T) (chi.Router, *services.MaterialService) {
	t.Helper()
	store := storage.New()
	materialSvc := services.NewMaterialService(store)
	handler := handlers.NewMaterialHandler(materialSvc)

	r := chi.NewRouter()
	r.Get("/material", handler.ListandoMateriales)
	r.Get("/material/{id}", handler.ObtenerMaterialPorID)
	r.Post("/material", handler.CreandoMaterial)
	r.Put("/material/{id}", handler.ActulizarUnMaterial)
	r.Delete("/material/{id}", handler.BorrarUnMaterial)

	return r, materialSvc
}

func TestMaterialHandler_ListarVacio(t *testing.T) {
	r, _ := setupMaterialRouter(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/material", nil)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	lista, ok := resp["data"].([]interface{})
	require.True(t, ok)
	assert.Empty(t, lista)
}

func TestMaterialHandler_CrearYObtener(t *testing.T) {
	r, _ := setupMaterialRouter(t)

	t.Run("crear valido -> 201", func(t *testing.T) {
		body := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"kg","precio_referencia":"25.50"}`
		req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "Cemento", data["nombre"])
	})

	t.Run("listar con 1 elemento", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/material", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		lista, ok := resp["data"].([]interface{})
		require.True(t, ok)
		assert.Len(t, lista, 1)
	})

	t.Run("obtener por id -> 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/material/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener id inexistente -> 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/material/999", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMaterialHandler_Errores(t *testing.T) {
	r, _ := setupMaterialRouter(t)

	t.Run("crear con nombre vacio -> 400", func(t *testing.T) {
		body := `{"nombre":"","descripcion":"Saco","unidad":"kg","precio_referencia":"25.50"}`
		req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("crear con JSON malformado -> 400", func(t *testing.T) {
		body := `{not-json}`
		req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("eliminar id inexistente -> 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/material/999", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("actualizar id inexistente -> 404", func(t *testing.T) {
		body := `{"nombre":"Nuevo","descripcion":"Item","unidad":"kg","precio_referencia":"10.00"}`
		req := httptest.NewRequest(http.MethodPut, "/material/999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMaterialHandler_ActualizarYEliminar(t *testing.T) {
	r, materialSvc := setupMaterialRouter(t)

	_, err := materialSvc.CrearM(models.EntradaMaterial{
		Nombre:           "Test",
		Descripcion:      "Item de prueba",
		Unidad:           "kg",
		PrecioReferencia: "10.00",
	})
	require.NoError(t, err)

	t.Run("actualizar -> 200", func(t *testing.T) {
		body := `{"nombre":"Actualizado","descripcion":"Item actualizado","unidad":"kg","precio_referencia":"15.00"}`
		req := httptest.NewRequest(http.MethodPut, "/material/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/material/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("verificar eliminado -> 404", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/material/1", nil)
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}