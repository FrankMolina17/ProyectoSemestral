package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)	


//Respuestas JSON

func respondJSON(w http.ResponseWriter, status int, payload any) { //esto es para manejar los errores
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
func respondOK(w http.ResponseWriter, data any) { //esto es para obtener un equipo
	respondJSON(w, http.StatusOK, map[string]any{"data": data})
}

func respondCreated(w http.ResponseWriter, data any, id int) { //esto es para crear un equipo
	respondJSON(w, http.StatusCreated, map[string]any{"data": data, "id": id})
}

func respondError(w http.ResponseWriter, status int, msg string) { //esto es para manejar los errores
	respondJSON(w, status, map[string]string{"error": msg})
}

//Manejo de Errores DE MEMORIA

func mapStoreError(w http.ResponseWriter, err error, recurso string, id int) { //esto es para manejar los errores
	switch {
	case errors.Is(err, storage.ErrNotFound):
		respondError(w, http.StatusNotFound,
			fmt.Sprintf("%s con id %d no encontrado", recurso, id))
	case errors.Is(err, storage.ErrDuplicated):
		respondError(w, http.StatusConflict,
			"nombre ya existe para esa unidad")
	default:
		respondError(w, http.StatusInternalServerError, "error interno")
	}
}


func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool { //esto es para decodificar el json
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		respondError(w, http.StatusBadRequest, "body malformado: "+err.Error())
		return false
	}
	return true
}

func urlParamID(w http.ResponseWriter, r *http.Request, param string) (int, bool) { //esto es para obtener el id
	raw := chi.URLParam(r, param)
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "id debe ser un entero positivo")
		return 0, false
	}
	return id, true
}
