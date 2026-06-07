package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// ProformaHandler expone los endpoints HTTP del módulo de proformas.
type ProformaHandler struct {
	storage *storage.ProformaStorage
}

// NewProformaHandler crea un handler con el almacenamiento indicado.
func NewProformaHandler(s *storage.ProformaStorage) *ProformaHandler {
	return &ProformaHandler{storage: s}
}

// CrearProforma registra una nueva proforma.
func (h *ProformaHandler) CrearProforma(w http.ResponseWriter, r *http.Request) {
	var proforma models.Proforma
	if err := json.NewDecoder(r.Body).Decode(&proforma); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo de solicitud inválido")
		return
	}

	creada, err := h.storage.Crear(&proforma)
	if err != nil {
		responderErrorStorage(w, err)
		return
	}

	responderJSON(w, http.StatusCreated, creada)
}

// ListarProformas devuelve todas las proformas.
func (h *ProformaHandler) ListarProformas(w http.ResponseWriter, r *http.Request) {
	proformas := h.storage.Listar()
	responderJSON(w, http.StatusOK, proformas)
}

// ObtenerProforma devuelve una proforma por ID.
func (h *ProformaHandler) ObtenerProforma(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	proforma, err := h.storage.ObtenerPorID(id)
	if err != nil {
		responderErrorStorage(w, err)
		return
	}

	responderJSON(w, http.StatusOK, proforma)
}

// ActualizarProforma modifica una proforma existente.
func (h *ProformaHandler) ActualizarProforma(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var proforma models.Proforma
	if err := json.NewDecoder(r.Body).Decode(&proforma); err != nil {
		responderError(w, http.StatusBadRequest, "cuerpo de solicitud inválido")
		return
	}

	actualizada, err := h.storage.Actualizar(id, proforma)
	if err != nil {
		responderErrorStorage(w, err)
		return
	}

	responderJSON(w, http.StatusOK, actualizada)
}

// CalcularProforma recalcula los costos de una proforma.
func (h *ProformaHandler) CalcularProforma(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		responderError(w, http.StatusBadRequest, "id inválido")
		return
	}

	proforma, err := h.storage.Calcular(id)
	if err != nil {
		responderErrorStorage(w, err)
		return
	}

	responderJSON(w, http.StatusOK, proforma)
}

func parseID(raw string) (int, error) {
	return strconv.Atoi(raw)
}

func responderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{"error": mensaje})
}

func responderErrorStorage(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrProformaNoEncontrada):
		responderError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, storage.ErrNombreRequerido):
		responderError(w, http.StatusBadRequest, err.Error())
	default:
		responderError(w, http.StatusInternalServerError, "error interno del servidor")
	}
}
