package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthRouter() (chi.Router, *services.AutenticacionService, *storage.Storage) {
	fakeStorage := storage.New()
	authSvc := services.NuevaAutenticacionService(fakeStorage, services.AuthOptions{})
	srv := handlers.NewServerC(nil, nil, nil, nil, authSvc)

	r := chi.NewRouter()
	r.Post("/api/v1/auth/register", srv.RegistrarUser)
	r.Post("/api/v1/auth/login", srv.LoginUser)
	r.Get("/api/v1/auth/usuarios", srv.ListarUsuarios)
	r.Get("/api/v1/auth/usuarios/{id}", srv.ObtenerUsuarioPorID)

	return r, authSvc, fakeStorage
}

func jsonReq(method, path, body, token string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func ejecutar(r chi.Router, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestRegister(t *testing.T) {
	h, _, _ := setupAuthRouter()

	t.Run("valido -> 201", func(t *testing.T) {
		body := `{"email":"nuevo@uleam.edu.ec","password":"secreta123"}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, ""))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
	t.Run("email duplicado -> 409", func(t *testing.T) {
		body := `{"email":"dup@uleam.edu.ec","password":"secreta123"}`
		ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, "")) // primero
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", body, "")) // repetido
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
	t.Run("JSON malformado -> 400", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", `{roto`, ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestLogin(t *testing.T) {
	h, _, _ := setupAuthRouter()
	cred := `{"email":"ana@uleam.edu.ec","password":"secreta123"}`
	ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", cred, ""))

	t.Run("correcto -> 200 + token", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/login", cred, ""))
		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Token string `json:"token"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		assert.NotEmpty(t, resp.Token)
	})
	t.Run("contrasena incorrecta -> 401", func(t *testing.T) {
		malo := `{"email":"ana@uleam.edu.ec","password":"incorrecta"}`
		rec := ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/login", malo, ""))
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestListarUsuarios(t *testing.T) {
	h, _, _ := setupAuthRouter()

	// Register 2 users
	ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", `{"email":"u1@test.com","password":"123456"}`, ""))
	ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", `{"email":"u2@test.com","password":"123456"}`, ""))

	rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/auth/usuarios", "", ""))
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	data, ok := resp["data"].([]interface{})
	require.True(t, ok)
	assert.Len(t, data, 2)
}

func TestObtenerUsuarioPorID(t *testing.T) {
	h, _, _ := setupAuthRouter()

	ejecutar(h, jsonReq(http.MethodPost, "/api/v1/auth/register", `{"email":"user@test.com","password":"123456"}`, ""))

	t.Run("existente -> 200", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/auth/usuarios/1", "", ""))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("inexistente -> 404", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/auth/usuarios/999", "", ""))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("id invalido -> 400", func(t *testing.T) {
		rec := ejecutar(h, jsonReq(http.MethodGet, "/api/v1/auth/usuarios/abc", "", ""))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
