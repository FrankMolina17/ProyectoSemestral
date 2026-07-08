package handlers

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// RespondJSON escribe data como JSON con el status HTTP indicado.
//
// Centraliza tres cosas que antes repetíamos en CADA handler:
//   - poner el header Content-Type
//   - escribir el status code
//   - codificar el cuerpo y registrar el error si la codificación falla
//
// Si data es nil (por ejemplo en un 204 No Content) no escribe cuerpo.
func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error codificando JSON: %v", err)
	}
}

// RespondError escribe un error en un formato JSON consistente: {"error": "..."}.
//
// Así el cliente siempre recibe los errores con la misma forma, en lugar de
// texto plano unas veces y JSON otras.
func RespondError(w http.ResponseWriter, status int, mensaje string) {
	RespondJSON(w, status, map[string]string{"error": mensaje})
}

func StatusDeError(err error) int {
	switch {
	case errors.Is(err, services.ErrNoEncontrado):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrCredencialesInvalidas):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrEmailEnUso):
		return http.StatusConflict
	case errors.Is(err, services.ErrNombreVacio):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrPrecioNegativo):
		return http.StatusConflict
	case errors.Is(err, services.ErrCredencialesInvalidos):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}

}
