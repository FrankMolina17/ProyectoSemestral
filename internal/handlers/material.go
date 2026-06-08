package handlers

import (
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// MaterialHandler maneja las operaciones relacionadas con los materiales.
type MaterialHandler struct {
	s *storage.Storage
}

// NewMaterialHandler crea un nuevo manejador de materiales.
func NewMaterialHandler(s *storage.Storage) *MaterialHandler {
	return &MaterialHandler{s: s}
}


// GET /api/v1/catalogo/material
func (h *MaterialHandler) Lista(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	nombre := q.Get("nombre")
	unidad := q.Get("unidad")

	list := h.s.ListarMateriales(nombre, unidad)
	respondOK(w, list)
}

// GET 
// Obtiene un material por su ID.
func (h *MaterialHandler) ObtenerUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	mat, err := h.s.ObtenerMaterialporID(id)
	if err != nil {
		mapStoreError(w, err, "material", id)
		return
	}
	respondOK(w, mat)
}

// POST /materiales
func (h *MaterialHandler) CrearUnMaterial(w http.ResponseWriter, r *http.Request) {
	var in models.MaterialInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	mat, err := h.s.CrearMaterial(in)
	if err != nil {
		mapStoreError(w, err, "material", 0)
		return
	}
	respondCreated(w, mat, mat.ID)
}

// PUT 
// Actualiza un material existente.
func (h *MaterialHandler) ActulizarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	var in models.MaterialInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	mat, err := h.s.ActualizarMaterial(id, in)
	if err != nil {
		mapStoreError(w, err, "material", id)
		return
	}
	respondOK(w, mat)
}

// DELETE 
// Elimina un material existente.
func (h *MaterialHandler) BorrarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	if err := h.s.BorrarMaterial(id); err != nil {
		mapStoreError(w, err, "material", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

