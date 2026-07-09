package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRepuestaJSON(t *testing.T) {
	w := httptest.NewRecorder()
	RepuestaJSON(w, http.StatusCreated, map[string]string{"msg": "ok"})
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestOk(t *testing.T) {
	w := httptest.NewRecorder()
	Ok(w, map[string]string{"msg": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreando(t *testing.T) {
	w := httptest.NewRecorder()
	Creando(w, map[string]string{"msg": "creado"}, 42)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestMalFormado(t *testing.T) {
	w := httptest.NewRecorder()
	MalFormado(w, "error de validacion")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNoEncontrado(t *testing.T) {
	w := httptest.NewRecorder()
	NoEncontrado(w, "material", 42)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestErrorMermoria(t *testing.T) {
	t.Run("ErrNotFound -> 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		ErrorMermoria(w, storage.ErrNotFound, "material", 1)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ErrDuplicated -> 409", func(t *testing.T) {
		w := httptest.NewRecorder()
		ErrorMermoria(w, storage.ErrDuplicated, "material", 1)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("other error -> 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		ErrorMermoria(w, errors.New("otro error"), "material", 1)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDecodificarJSON(t *testing.T) {
	t.Run("JSON valido -> true", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test"}`))
		r.Header.Set("Content-Type", "application/json")
		var target struct {
			Name string `json:"name"`
		}
		valid := DecodificarJSON(w, r, &target)
		assert.True(t, valid)
		assert.Equal(t, "test", target.Name)
	})

	t.Run("body malformado -> false", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{roto}`))
		r.Header.Set("Content-Type", "application/json")
		var target struct{}
		valid := DecodificarJSON(w, r, &target)
		assert.False(t, valid)
	})
}

func setURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		rctx = chi.NewRouteContext()
	}
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestParaObtenerelID(t *testing.T) {
	t.Run("id valido", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := setURLParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "5")

		id, valid := ParaObtenerelID(w, r)
		assert.True(t, valid)
		assert.Equal(t, 5, id)
	})

	t.Run("id invalido -> false", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := setURLParam(httptest.NewRequest(http.MethodGet, "/", nil), "id", "abc")

		_, valid := ParaObtenerelID(w, r)
		assert.False(t, valid)
	})
}

func TestParaObtenerTipoRecursoID(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := setURLParam(httptest.NewRequest(http.MethodGet, "/", nil), "tipo", "material")
		r = setURLParam(r, "recursoID", "3")

		tipo, id, valid := ParaObtenerTipoRecursoID(w, r)
		assert.True(t, valid)
		assert.Equal(t, "material", tipo)
		assert.Equal(t, 3, id)
	})

	t.Run("tipo invalido -> false", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := setURLParam(httptest.NewRequest(http.MethodGet, "/", nil), "tipo", "invalido")
		r = setURLParam(r, "recursoID", "1")

		_, _, valid := ParaObtenerTipoRecursoID(w, r)
		assert.False(t, valid)
	})

	t.Run("recursoID invalido -> false", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := setURLParam(httptest.NewRequest(http.MethodGet, "/", nil), "tipo", "material")
		r = setURLParam(r, "recursoID", "abc")

		_, _, valid := ParaObtenerTipoRecursoID(w, r)
		assert.False(t, valid)
	})
}
