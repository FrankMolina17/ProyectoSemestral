package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func obtenerTokenValido(t *testing.T, authSvc *services.AutenticacionService, s *storage.Storage) string {
	t.Helper()
	u, err := s.CrearUsuario(models.EntradaUsuario{Email: "test@test.com", Password: "123456"})
	require.NoError(t, err)
	token, err := authSvc.GenerarJWT(*u)
	require.NoError(t, err)
	return token
}

func setupManoObraRouter() (chi.Router, *services.AutenticacionService, *storage.Storage) {
	fakeStorage := storage.New()
	authSvc := services.NuevaAutenticacionService(fakeStorage, services.AuthOptions{})
	mh := NewManoObraHandler(services.NewManoObraService(fakeStorage))

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthJWT(authSvc))
		r.Get("/manoobra", mh.ListarUnaManoObra)
		r.Get("/manoobra/{id}", mh.ObtenerUnaManoObraPorID)
		r.Post("/manoobra", mh.CreandoUnaManoObra)
		r.Put("/manoobra/{id}", mh.ActualizadoUnaManoObra)
		r.Delete("/manoobra/{id}", mh.BorrandoUnaManoObra)
	})

	return r, authSvc, fakeStorage
}

func TestListarManoObra_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupManoObraRouter()

	req := httptest.NewRequest(http.MethodGet, "/manoobra", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestObtenerManoObraPorID_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupManoObraRouter()

	req := httptest.NewRequest(http.MethodGet, "/manoobra/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestCrearManoObra_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupManoObraRouter()

	body := `{"descripcion":"Oficial","categoria":"oficial","unidad":"día","costo_referencia":10.00}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestCrearManoObra_ConTokenYDescripcionVacia_Devuelve400(t *testing.T) {
	r, authSvc, s := setupManoObraRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"descripcion":"","categoria":"oficial","unidad":"día","costo_referencia":10.00}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "descripcion")
}

func TestCrearManoObra_ConTokenYCategoriaInvalida_Devuelve400(t *testing.T) {
	r, authSvc, s := setupManoObraRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"descripcion":"Oficial","categoria":"invalida","unidad":"día","costo_referencia":10.00}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "categoria")
}

func TestCrearManoObra_ConTokenYUnidadInvalida_Devuelve400(t *testing.T) {
	r, authSvc, s := setupManoObraRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"descripcion":"Oficial","categoria":"oficial","unidad":"km","costo_referencia":10.00}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "unidad")
}

func TestCrearManoObra_ConTokenYCostoCero_Devuelve400(t *testing.T) {
	r, authSvc, s := setupManoObraRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"descripcion":"Oficial","categoria":"oficial","unidad":"día","costo_referencia":0}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "costo")
}

func TestCrearManoObra_ConTokenYDatosValidos_Devuelve201YNoDuplica(t *testing.T) {
	r, authSvc, s := setupManoObraRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"descripcion":"Oficial albañil","categoria":"oficial","unidad":"día","costo_referencia":10.00}`
	req := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Oficial albañil", data["descripcion"])
	assert.Equal(t, "oficial", data["categoria"])
	assert.Equal(t, float64(2), resp["id"])

	body2 := `{"descripcion":"Oficial albañil","categoria":"oficial","unidad":"día","costo_referencia":10.00}`
	req2 := httptest.NewRequest(http.MethodPost, "/manoobra", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
	var resp2 map[string]string
	json.NewDecoder(w2.Body).Decode(&resp2)
	assert.Contains(t, resp2["error"], "ya existe")
}
