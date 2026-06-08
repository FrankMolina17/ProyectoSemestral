package handlers

import (
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

type ManoObraHandler struct{ s *storage.Storage } 

func NewManoObraHandler(s *storage.Storage) *ManoObraHandler { //mano de obra
	return &ManoObraHandler{s: s}
}

// GET api/v1/manoobra
func (h *ManoObraHandler) List(w http.ResponseWriter, r *http.Request) {
	categoria := r.URL.Query().Get("categoria")
	list := h.s.ListarManoObra(categoria)
	respondOK(w, list)
}

// GET /mano-obra/:id
func (h *ManoObraHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	mo, err := h.s.ObtenerManoObraPorID(id)
	if err != nil {
		mapStoreError(w, err, "mano_obra", id)
		return
	}
	respondOK(w, mo)
}

// POST /mano-obra
func (h *ManoObraHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.ManoObraInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	mo, err := h.s.CrearManoObra(in)
	if err != nil {
		mapStoreError(w, err, "mano_obra", 0)
		return
	}
	respondCreated(w, mo, mo.ID)
}

// PUT /mano-obra/:id
func (h *ManoObraHandler) Replace(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	var in models.ManoObraInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	mo, err := h.s.ActualizarManoObra(id, in)
	if err != nil {
		mapStoreError(w, err, "mano_obra", id)
		return
	}
	respondOK(w, mo)
}

// DELETE /mano-obra/:id
func (h *ManoObraHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	if err := h.s.BorrarManoObra(id); err != nil {
		mapStoreError(w, err, "mano_obra", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
