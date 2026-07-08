package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

// CrearObra atiende POST /api/v1/obras
func (s *Server) CrearObra(w http.ResponseWriter, r *http.Request) {
	var obra models.Obra
	if err := json.NewDecoder(r.Body).Decode(&obra); err != nil {
		RespondJSON(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	nueva, err := s.ObraService.CrearObra(obra)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, nueva)
}

// ObtenerObras atiende GET /api/v1/obras
func (s *Server) ObtenerObras(w http.ResponseWriter, r *http.Request) {
	obras := s.ObraService.Listar()
	RespondJSON(w, http.StatusOK, obras)
}

// ObtenerObraPorID atiende GET /api/v1/obras/{id}
func (s *Server) ObtenerObraPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	obra, err := s.ObraService.Obtener(id)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Obra no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, obra)
}

// ActualizarObra atiende PUT /api/v1/obras/{id}
func (s *Server) ActualizarObra(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	var datos models.Obra
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondJSON(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if strings.TrimSpace(datos.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "el campo nombre es obligatorio")
		return
	}

	actualizada, err := s.ObraService.ActualizarObra(id, datos)
	if err != nil {
		RespondError(w, http.StatusNotFound, "obra no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// EliminarObra atiende DELETE /api/v1/obras/{id}
func (s *Server) EliminarObra(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	if err := s.ObraService.BorrarObra(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
