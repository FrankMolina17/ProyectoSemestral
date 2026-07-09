package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/fakes"
	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() http.Handler {
	fakeRepo := fakes.NuevoProformaRepositoryFake()
	proformaSvc := services.NuevoProformaService(fakeRepo)
	proformaHandler := handlers.NuevoHandler(proformaSvc)

	usuarioStore := storage.NuevoUsuarioStorage()
	authService := services.NuevoAuthService(usuarioStore)

	r := chi.NewRouter()
	r.Use(middleware.CORS)

	authHandler := handlers.NuevoAuthHandler(authService)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Registrar)
		r.Post("/login", authHandler.Login)
		r.Post("/register-admin", authHandler.RegistrarAdmin)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.VerificarJWT(authService))

		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/proformas", proformaHandler.CrearProforma)
			r.Get("/proformas", proformaHandler.ObtenerTodos)
			r.Get("/proformas/{id}", proformaHandler.ObtenerPorID)
		})
	})

	return r
}

func obtenerToken(t *testing.T, router http.Handler) string {
	t.Helper()

	body := bytes.NewBufferString(`{"email":"test@obras.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	body = bytes.NewBufferString(`{"email":"test@obras.com","password":"secret123"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string 
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp["token"])
	return resp["token"]
}

// TestObtenerProformas_SinToken responde 401 cuando la ruta protegida
// se accede sin header Authorization.
func TestObtenerProformas_SinToken(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/proformas", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "token requerido")
}

// TestCrearProforma_ConToken usa un fake en memoria y verifica que
// una proforma válida se crea correctamente vía HTTP.
func TestCrearProforma_ConToken(t *testing.T) {
	router := setupTestRouter()
	token := obtenerToken(t, router)

	payload := `{"nombre":"Proforma remodelación","obra_id":7}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/proformas", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var proforma map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &proforma))
	assert.Equal(t, "Proforma remodelación", proforma["nombre"])
	assert.Equal(t, float64(7), proforma["obra_id"])
	assert.Equal(t, "borrador", proforma["estado"])
}

func TestRegistrarAdmin(t *testing.T) {
	router := setupTestRouter()

	t.Run("crear admin -> 201", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email":"admin@uleam.edu.ec","password":"admin123","rol":"admin"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register-admin", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "admin@uleam.edu.ec", resp["email"])
		assert.Equal(t, "admin", resp["rol"])
	})

	t.Run("sin rol -> default admin", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email":"admin2@uleam.edu.ec","password":"admin123"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register-admin", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "admin", resp["rol"])
	})

	t.Run("email duplicado -> 400", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email":"admin@uleam.edu.ec","password":"admin123","rol":"admin"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register-admin", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_Errores(t *testing.T) {
	router := setupTestRouter()

	t.Run("register JSON malformado -> 400", func(t *testing.T) {
		body := bytes.NewBufferString(`{roto}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("register email vacio -> 400", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email":"","password":"123456"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("login credenciales invalidas -> 401", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email":"noexiste@test.com","password":"wrong"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
