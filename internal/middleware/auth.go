// Package middleware provee funciones HTTP reutilizables.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthJWT valida el header Authorization: Bearer <token>.
// En producción se verifica la firma con una clave secreta (p.ej. golang-jwt/jwt).
// Aquí se acepta cualquier token no vacío para mantener el foco en el diseño de API.
func AuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "token inválido",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "token inválido",
			})
			return
		}

		// verificar firma JWT con golang-jwt/jwt y extraer claims reales.
		
		// Por ahora almacenamos el token crudo como "userID" en el contexto.
		ctx := context.WithValue(r.Context(), UserIDKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
