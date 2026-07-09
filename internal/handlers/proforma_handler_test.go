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

	r.Route("/api/v1/auth", func(r chi.Router) {
		authHandler := handlers.NuevoAuthHandler(authService)
		r.Post("/register", authHandler.Registrar)
		r.Post("/login", authHandler.Login)
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
