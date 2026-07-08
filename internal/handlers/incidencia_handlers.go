package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

// CrearIncidencia atiende POST /api/v1/incidencias
func (s *Server) CrearIncidencia(w http.ResponseWriter, r *http.Request) {
	var incidencia models.Incidencia
	if err := json.NewDecoder(r.Body).Decode(&incidencia); err != nil {
		RespondJSON(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	nueva, err := s.IncidenciaService.CrearIncidencia(incidencia)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, nueva)
}

// ObtenerIncidencias atiende GET /api/v1/incidencias
func (s *Server) ObtenerIncidencias(w http.ResponseWriter, r *http.Request) {
	incidencias := s.IncidenciaService.Listar()
	RespondJSON(w, http.StatusOK, incidencias)
}

// ObtenerIncidenciaPorID atiende GET /api/v1/incidencias/{id}
func (s *Server) ObtenerIncidenciaPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	incidencia, err := s.IncidenciaService.Obtener(id)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Incidencia no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, incidencia)
}

// ObtenerIncidenciasPorEntidad atiende GET /api/v1/incidencias/por/{tipo}/{id}
func (s *Server) ObtenerIncidenciasPorEntidad(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	tipo := chi.URLParam(r, "tipo")

	if tipo == "" {
		RespondError(w, http.StatusBadRequest, "Tipo es requerido")
		return
	}

	incidencia, err := s.IncidenciaService.ObtenerPorEntidad(id, tipo)

	if err != nil {
		RespondError(w, http.StatusNotFound, "Incidencia no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, incidencia)
}

// EliminarIncidencia atiende DELETE /api/v1/incidencias/{id}
func (s *Server) EliminarIncidencia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	if err := s.IncidenciaService.BorrarIncidencia(id); err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}

func (s *Server) ActualizarIncidencia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "id debe ser un número entero")
		return
	}

	var datos models.Incidencia
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if strings.TrimSpace(datos.Titulo) == "" {
		RespondError(w, http.StatusBadRequest, "el campo nombre es obligatorio")
		return
	}

	actualizada, err := s.IncidenciaService.ActualizarIncidencia(id, datos)

	if err != nil {
		RespondError(w, http.StatusNotFound, "categoría no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}
