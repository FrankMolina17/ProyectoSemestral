package handlers

import (
	"net/http"
	"strconv"


	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

type EquipoHandler struct{ s *storage.Storage }

func NewEquipoHandler(s *storage.Storage) *EquipoHandler {
	return &EquipoHandler{s: s}
}

//ListarEquipo
func (h *EquipoHandler) Lista(w http.ResponseWriter, r *http.Request) { //esto es para obtener todos los equipos
	q := r.URL.Query() 
	tipo := q.Get("tipo")

	var disponible *bool
	if raw := q.Get("disponible"); raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			respondError(w, http.StatusBadRequest, "disponible debe ser true o false")
			return
		}
		disponible = &b
	}

	list := h.s.ListarEquipos(disponible, tipo)
	respondOK(w, list)
}

func (h *EquipoHandler) GetByID(w http.ResponseWriter, r *http.Request) { //esto es para obtener un equipo
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	eq, err := h.s.ObtenerEquipoPorID(id)
	if err != nil {
		mapStoreError(w, err, "equipo", id)
		return
	}
	respondOK(w, eq)
}

//Post
func (h *EquipoHandler) Crear(w http.ResponseWriter, r *http.Request) {//esto es para crear un equipo
	var in models.EquipoInput
	if !decodeJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		respondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	eq, err := h.s.CrearEquipo(in)
	if err != nil {
		mapStoreError(w, err, "equipo", 0)
		return
	}
	respondCreated(w, eq, eq.ID)
}

func (h *EquipoHandler) PatchDisponibilidad(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	var in models.DisponibilidadInput
	if !decodeJSON(w, r, &in) {
		return
	}
	eq, err := h.s.Disponibilidad(id, in.Disponible)
	if err != nil {
		mapStoreError(w, err, "equipo", id)
		return
	}
	respondOK(w, eq)
}

// DELETE /equipos/:id
func (h *EquipoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := urlParamID(w, r, "id")
	if !ok {
		return
	}
	if err := h.s.BorrarEquipo(id); err != nil {
		mapStoreError(w, err, "equipo", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
