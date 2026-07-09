package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupMaterialRouter() (chi.Router, *services.AutenticacionService, *storage.Storage) {
	fakeStorage := storage.New()
	authSvc := services.NuevaAutenticacionService(fakeStorage, services.AuthOptions{})
	mh := NewMaterialHandler(services.NewMaterialService(fakeStorage))

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthJWT(authSvc))
		r.Get("/material", mh.ListandoMateriales)
		r.Get("/material/{id}", mh.ObtenerMaterialPorID)
		r.Post("/material", mh.CreandoMaterial)
		r.Put("/material/{id}", mh.ActulizarUnMaterial)
		r.Delete("/material/{id}", mh.BorrarUnMaterial)
	})

	return r, authSvc, fakeStorage
}

func obtenerTokenValido(t *testing.T, authSvc *services.AutenticacionService, s *storage.Storage) string {
	_, err := authSvc.RegistrarUsuario("test@example.com", "password123")
	assert.NoError(t, err)

	usuario, err := authSvc.Login("test@example.com", "password123")
	assert.NoError(t, err)

	token, err := authSvc.GenerarJWT(*usuario)
	assert.NoError(t, err)
	return token
}

func TestCrearMaterial_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupMaterialRouter()

	body := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"unidad","precio_referencia":"25.50"}`
	req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestListarMateriales_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupMaterialRouter()

	req := httptest.NewRequest(http.MethodGet, "/material", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestObtenerMaterialPorID_SinToken_Devuelve401(t *testing.T) {
	r, _, _ := setupMaterialRouter()

	req := httptest.NewRequest(http.MethodGet, "/material/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autorizado")
}

func TestCrearMaterial_ConTokenYNombreVacio_Devuelve400(t *testing.T) {
	r, authSvc, s := setupMaterialRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"nombre":"","descripcion":"Saco 50kg","unidad":"unidad","precio_referencia":"25.50"}`
	req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "nombre")
}

func TestCrearMaterial_ConTokenYUnidadInvalida_Devuelve400(t *testing.T) {
	r, authSvc, s := setupMaterialRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"km","precio_referencia":"25.50"}`
	req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "unidad")
}

func TestCrearMaterial_ConTokenYPrecioCero_Devuelve400(t *testing.T) {
	r, authSvc, s := setupMaterialRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"unidad","precio_referencia":"0"}`
	req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Contains(t, resp["error"], "precio")
}

func TestCrearMaterial_ConTokenYDatosValidos_Devuelve201YNoDuplica(t *testing.T) {
	r, authSvc, s := setupMaterialRouter()
	token := obtenerTokenValido(t, authSvc, s)

	body := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"unidad","precio_referencia":"25.50"}`
	req := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	data := resp["data"].(map[string]any)
	assert.Equal(t, "Cemento", data["nombre"])
	assert.Equal(t, "unidad", data["unidad"])
	assert.Equal(t, float64(2), resp["id"])

	body2 := `{"nombre":"Cemento","descripcion":"Saco 50kg","unidad":"unidad","precio_referencia":"25.50"}`
	req2 := httptest.NewRequest(http.MethodPost, "/material", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)
	var resp2 map[string]string
	json.NewDecoder(w2.Body).Decode(&resp2)
	assert.Contains(t, resp2["error"], "nombre ya existe")
}
