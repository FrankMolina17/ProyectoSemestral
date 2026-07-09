package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"

	"github.com/stretchr/testify/assert"
)

func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["msg"])
}

func TestRespondJSON_Nil(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondJSON(w, http.StatusNoContent, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondError(w, http.StatusBadRequest, "algo salio mal")
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "algo salio mal", resp["error"])
}

func TestStatusDeError(t *testing.T) {
	w := httptest.NewRecorder()

	// Probamos el helper de status de error indirectamente
	handlers.RespondJSON(w, http.StatusConflict, map[string]string{"error": "conflicto"})
	assert.Equal(t, http.StatusConflict, w.Code)
}
