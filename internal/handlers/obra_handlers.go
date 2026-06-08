package handlers

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func CrearObraHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var obra models.Obra
	if err := json.NewDecoder(r.Body).Decode(&obra); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	nueva, err := services.CrearObra(obra)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nueva)
}

func ObtenerObrasHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	obras := services.ObtenerObras()
	json.NewEncoder(w).Encode(obras)
}

func ObtenerObraHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	obra, encontrado := services.ObtenerObra(id)
	if !encontrado {
		http.Error(w, "Obra no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(obra)
}

func ActualizarObraHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var obra models.Obra
	if err := json.NewDecoder(r.Body).Decode(&obra); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	obra.ID = id
	obraActualizada, actualizado := services.ActualizarObra(id, obra)
	if !actualizado {
		http.Error(w, "No se pudo actualizar la obra", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(obraActualizada)
}

func EliminarObraHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	if !services.EliminarObra(id) {
		http.Error(w, "No se pudo eliminar la obra", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
