package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"


	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// PrecioHandler maneja las operaciones relacionadas con los precios de recursos.
type PrecioHandler struct {
	s *storage.Storage
}

// NewPrecioHandler crea un nuevo manejador de precios.
func NewPrecioHandler(s *storage.Storage) *PrecioHandler {
	return &PrecioHandler{s: s}
}

// GET /precios/:tipo/:id
// tipo ∈ {material, mano_obra, equipo}
//historial de precios de un recurso
func (h *PrecioHandler) Historial(w http.ResponseWriter, r *http.Request) {
	tipo, id, ok := h.parseTipoID(w, r)
	if !ok {
		return
	}
	list := h.s.HistorialPrecios(tipo, id)
	respondOK(w, list)
}

// GET /precios/:tipo/:id/vigente
func (h *PrecioHandler) Vigente(w http.ResponseWriter, r *http.Request) {
	tipo, id, ok := h.parseTipoID(w, r)
	if !ok {
		return
	}
	p, err := h.s.PrecioVigente(tipo, id)
	if err != nil {
		respondError(w, http.StatusNotFound,
			"no hay precio vigente registrado para ese recurso")
		return
	}
	respondOK(w, p)
}

// POST /precios
func (h *PrecioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.PrecioRecursoInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	p, err := h.s.CreatePrecio(in)
	if err != nil {
		mapStoreError(w, err, in.RecursoTipo, in.RecursoID)
		return
	}
	respondCreated(w, p, p.ID)
}

// parseTipoID parsea el tipo de recurso y su ID de la URL.
func (h *PrecioHandler) parseTipoID(w http.ResponseWriter, r *http.Request) (string, int, bool) {
	tipo := chi.URLParam(r, "tipo")
	if !models.RecursosTipos[tipo] {
		respondError(w, http.StatusBadRequest,
			"tipo no válido: use material, mano_obra, equipo")
		return "", 0, false
	}
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return "", 0, false
	}
	return tipo, id, true
}
