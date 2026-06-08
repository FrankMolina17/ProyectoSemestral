package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/go-chi/chi/v5"
)

func CrearIncidenciaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var incidencia models.Incidencia
	if err := json.NewDecoder(r.Body).Decode(&incidencia); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	nueva, err := services.CrearIncidencia(incidencia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nueva)
}

func ObtenerIncidenciasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	incidencias := services.ObtenerIncidencias()
	json.NewEncoder(w).Encode(incidencias)
}

func ObtenerIncidenciaPorIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	incidencia, encontrado := services.ObtenerIncidenciaPorID(id)
	if !encontrado {
		http.Error(w, "Incidencia no encontrada", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(incidencia)
}

func ObtenerIncidenciasPorEntidadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	entidadTipo := chi.URLParam(r, "tipo")
	idStr := chi.URLParam(r, "id")

	entidadID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	incidencias := services.ObtenerIncidenciasPorEntidad(entidadTipo, entidadID)
	json.NewEncoder(w).Encode(incidencias)
}

func ActualizarIncidenciaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var datos models.Incidencia
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	actualizada, err := services.ActualizarIncidencia(id, datos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(actualizada)
}

func EliminarIncidenciaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = services.EliminarIncidencia(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
