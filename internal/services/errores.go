package services

import (
	"errors"

	"encoding/json"

	"fmt"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
)

var (
	ErrEmailVacio                = errors.New("email y password son requeridos")
	ErrRecursoNoEncontrado       = errors.New("recurso no encontrado")
	ErrNombreVacio               = errors.New("el campo es obligatorio")
	ErrUnidadVacia               = errors.New("el campo unidad es obligatorio")
	ErrDescripcionVacia          = errors.New("el campo descripcion es obligatorio")
	ErrUnidadNoPermitida         = errors.New("unidad no permitida")
	ErrPrecioReferencialInvalido = errors.New("debe ser mayor a 0")
	ErrNoEncontrado              = errors.New("registro no encontrado")
	ErrPrecioNegativo            = errors.New("el precio no puede ser negativo")
	ErrEmailEnUso                = errors.New("el correo electrónica ya se encuentra en uso")
	ErrCredencialesInvalidas     = errors.New("Email o contraseña incorrectos")
	ErrNotFound                  = errors.New("recurso no encontrado")
	ErrDuplicated                = errors.New("nombre ya existe para esa unidad")
	ErrCategoriaNoPermitida      = errors.New("categoria no permitida")
	ErrTipoVacio                 = errors.New("el campo tipo es obligatorio")
	ErrFechaVigenciaVacia        = errors.New("el campo fecha vigencia es obligatorio")
)

func RepuestaJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// su funcion es devolver un 200 en caso de exito
func Ok(w http.ResponseWriter, data any) {
	RepuestaJSON(w, http.StatusOK, map[string]any{"data": data})
}
func Creando(w http.ResponseWriter, data any, id int) {
	RepuestaJSON(w, http.StatusCreated, map[string]any{"data": data, "id": id})
}

// su funcion es devolver un 400
func MalFormado(w http.ResponseWriter, msg string) {
	RepuestaJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}

// su funcion es devolver un 404 en caso de no encontrar el recurso
func NoEncontrado(w http.ResponseWriter, recurso string, id int) {
	RepuestaJSON(w, http.StatusNotFound,
		map[string]string{"error": fmt.Sprintf("%s con id %d no encontrado", recurso, id)})
}

// su funcion es devolver un 500 en caso de error de servidor
func ErrorMermoria(w http.ResponseWriter, err error, recurso string, id int) {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		NoEncontrado(w, recurso, id)
	case errors.Is(err, storage.ErrDuplicated):
		RepuestaJSON(w, http.StatusConflict,
			map[string]string{"error": "nombre ya existe para esa unidad"})
	case err != nil:
		RepuestaJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "error interno"})
	}
}

// esto es para decodificar el json y mostrar el error
func DecodificarJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		MalFormado(w, "body malformado: "+err.Error())
		return false
	}
	return true
}

// esto es para obtener el id
func ParaObtenerelID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		MalFormado(w, "id debe ser un entero positivo")
		return 0, false
	}
	return id, true
}

func ParaObtenerTipoRecursoID(w http.ResponseWriter, r *http.Request) (string, int, bool) {
	tipo := chi.URLParam(r, "tipo")
	if !models.RecursosTipos[tipo] {
		MalFormado(w, "recurso_tipo debe ser: material, mano_obra, equipo")
		return "", 0, false
	}
	recursoID, err := strconv.Atoi(chi.URLParam(r, "recursoID"))
	if err != nil || recursoID <= 0 {
		MalFormado(w, "recurso_id debe ser un entero positivo")
		return "", 0, false
	}
	return tipo, recursoID, true
}
