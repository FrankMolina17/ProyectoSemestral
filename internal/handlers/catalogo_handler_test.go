package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
// ManoObra
// ─────────────────────────────────────────────

func newManoObraRouter(t *testing.T) (chi.Router, *storage.Storage) {
	t.Helper()
	store := storage.New()
	svc := services.NewManoObraService(store)
	h := handlers.NewManoObraHandler(svc)

	r := chi.NewRouter()
	r.Get("/manoobra", h.ListarUnaManoObra)
	r.Get("/manoobra/{id}", h.ObtenerUnaManoObraPorID)
	r.Post("/manoobra", h.CreandoUnaManoObra)
	r.Put("/manoobra/{id}", h.ActualizadoUnaManoObra)
	r.Delete("/manoobra/{id}", h.BorrandoUnaManoObra)
	return r, store
}

func TestManoObraHandler_FullCRUD(t *testing.T) {
	r, _ := newManoObraRouter(t)

	t.Run("listar vacio -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/manoobra", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener por id inexistente -> 404", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/manoobra/999", "", ""))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	var createdID int
	t.Run("crear -> 201", func(t *testing.T) {
		body := `{"descripcion":"Albañil","categoria":"oficial","unidad":"hora","costo_referencia":15.00}`
		rec := ejecutar(r, jsonReq(http.MethodPost, "/manoobra", body, ""))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		data, ok := resp["data"].(map[string]interface{})
		require.True(t, ok)
		createdID = int(data["id"].(float64))
	})

	t.Run("obtener por id existente -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, fmt.Sprintf("/manoobra/%d", createdID), "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("actualizar -> 200", func(t *testing.T) {
		body := `{"descripcion":"Albañil Senior","categoria":"oficial","unidad":"hora","costo_referencia":20.00}`
		rec := ejecutar(r, jsonReq(http.MethodPut, fmt.Sprintf("/manoobra/%d", createdID), body, ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodPost, "/manoobra", `{roto`, ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodDelete, fmt.Sprintf("/manoobra/%d", createdID), "", ""))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}

// ─────────────────────────────────────────────
// Equipo
// ─────────────────────────────────────────────

func newEquipoRouter(t *testing.T) (chi.Router, *storage.Storage) {
	t.Helper()
	store := storage.New()
	svc := services.NewEquipoService(store)
	h := handlers.NewEquipoHandler(svc)

	r := chi.NewRouter()
	r.Get("/equipos", h.ListandoLosEquipos)
	r.Get("/equipos/{id}", h.ObtenerUnEquipoPorID)
	r.Post("/equipos", h.CrearUnEquipo)
	r.Put("/equipos/{id}", h.ActualizarUnEquipo)
	r.Delete("/equipos/{id}", h.BorrarUnEquipo)
	return r, store
}

func TestEquipoHandler_FullCRUD(t *testing.T) {
	r, _ := newEquipoRouter(t)

	t.Run("listar vacio -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/equipos", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("obtener por id inexistente -> 404", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/equipos/999", "", ""))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	var createdID int
	t.Run("crear -> 201", func(t *testing.T) {
		body := `{"nombre":"Excavadora","tipo":"pesado","unidad":"hora","costo_hora":85.00}`
		rec := ejecutar(r, jsonReq(http.MethodPost, "/equipos", body, ""))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		data, ok := resp["data"].(map[string]interface{})
		require.True(t, ok)
		createdID = int(data["id"].(float64))
	})

	t.Run("obtener por id existente -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, fmt.Sprintf("/equipos/%d", createdID), "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("actualizar -> 200", func(t *testing.T) {
		body := `{"nombre":"Excavadora Grande","tipo":"pesado","unidad":"hora","costo_hora":95.00}`
		rec := ejecutar(r, jsonReq(http.MethodPut, fmt.Sprintf("/equipos/%d", createdID), body, ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodDelete, fmt.Sprintf("/equipos/%d", createdID), "", ""))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodPost, "/equipos", `{roto`, ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ─────────────────────────────────────────────
// Precio
// ─────────────────────────────────────────────

func newPrecioFullRouter(t *testing.T) (chi.Router, *storage.Storage) {
	t.Helper()
	store := storage.New()
	svc := services.NewPreciosService(store)
	h := handlers.NewPrecioHandler(svc)

	// Seed a material so there's a valid reference for precios
	store.CrearMateriales(models.EntradaMaterial{
		Nombre:           "Cemento",
		Unidad:           "unidad",
		PrecioReferencia: "10.00",
	})

	r := chi.NewRouter()
	r.Get("/precios", h.ListarDeLosPrecios)
	r.Get("/precios/{id}", h.ObtenerUnPrecioPorID)
	r.Post("/precios", h.CrearUnPrecio)
	r.Put("/precios/{id}", h.ActualizarUnPrecio)
	r.Delete("/precios/{id}", h.BorrarUnPrecio)
	r.Get("/precio/{tipo}/{recursoID}", h.HistorialPorRecurso)
	r.Get("/precio/{tipo}/{recursoID}/vigente", h.PrecioVigentePorRecurso)
	return r, store
}

func TestPrecioHandler_FullCRUD(t *testing.T) {
	r, store := newPrecioFullRouter(t)
	_ = store

	t.Run("listar vacio -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/precios", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("historial recurso inexistente -> 200 (vacio)", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/precio/material/999", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("vigente recurso inexistente -> 404", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/precio/material/999/vigente", "", ""))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	var createdID int
	t.Run("crear -> 201", func(t *testing.T) {
		fecha := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		body := fmt.Sprintf(`{"recurso_tipo":"material","recurso_id":1,"precio":12.50,"fecha_vigencia":"%s"}`, fecha)
		rec := ejecutar(r, jsonReq(http.MethodPost, "/precios", body, ""))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		createdID = int(resp["id"].(float64))
	})

	t.Run("obtener por id -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, fmt.Sprintf("/precios/%d", createdID), "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("historial -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/precio/material/1", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("vigente -> 200", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodGet, "/precio/material/1/vigente", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("actualizar -> 200", func(t *testing.T) {
		fecha := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		body := fmt.Sprintf(`{"recurso_tipo":"material","recurso_id":1,"precio":15.00,"fecha_vigencia":"%s"}`, fecha)
		rec := ejecutar(r, jsonReq(http.MethodPut, fmt.Sprintf("/precios/%d", createdID), body, ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("eliminar -> 204", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodDelete, fmt.Sprintf("/precios/%d", createdID), "", ""))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(r, jsonReq(http.MethodPost, "/precios", `{roto`, ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
