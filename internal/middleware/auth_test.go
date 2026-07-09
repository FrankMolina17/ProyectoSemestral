package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupAuth() (*services.AutenticacionService, *storage.Storage) {
	s := storage.New()
	auth := services.NuevaAutenticacionService(s, services.AuthOptions{})
	return auth, s
}

func mustToken(auth *services.AutenticacionService, usuario models.Usuario) string {
	t, err := auth.GenerarJWT(usuario)
	if err != nil {
		panic(err)
	}
	return t
}

func TestAuthJWT(t *testing.T) {
	auth, _ := setupAuth()
	token := mustToken(auth, models.Usuario{ID: 1, Email: "test@test.com", Rol: "admin"})

	r := chi.NewRouter()
	r.Use(middleware.AuthJWT(auth))
	r.Get("/protegido", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	})

	t.Run("token valido -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("sin token -> 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("token invalido -> 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protegido", nil)
		req.Header.Set("Authorization", "Bearer token-malo")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestRequerirRol(t *testing.T) {
	auth, _ := setupAuth()
	adminToken := mustToken(auth, models.Usuario{ID: 1, Email: "admin@test.com", Rol: "admin"})

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthJWT(auth))
		r.With(middleware.RequerirRol("admin")).Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		})
		r.With(middleware.RequerirRol("cliente")).Get("/cliente", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		})
	})

	t.Run("admin accede a /admin -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("admin accede a /cliente -> 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cliente", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
