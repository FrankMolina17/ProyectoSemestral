package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

// ==================== Auth para Catálogo (Módulo 1) ====================
type contextoClave string

const contextoClaveKey = contextoClave("usuarioID")

func AuthJWT(authSvc *services.AutenticacionService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				NoAutorizado(w)
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			token = strings.TrimSpace(token)

			claims, err := authSvc.ValidarJWT(token)
			if err != nil {
				NoAutorizado(w)
				return
			}

			ctx := context.WithValue(r.Context(), contextoClaveKey, claims.UsuarioID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ==================== Auth para Proformas (Módulo 2) ====================
func VerificarJWT(authSvc *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"token requerido"}`, http.StatusUnauthorized)
				return
			}

			partes := strings.SplitN(authHeader, " ", 2)
			if len(partes) != 2 || partes[0] != "Bearer" {
				http.Error(w, `{"error":"formato de token inválido"}`, http.StatusUnauthorized)
				return
			}

			_, err := authSvc.VerificarToken(partes[1])
			if err != nil {
				http.Error(w, `{"error":"token inválido o expirado"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ==================== Helper ====================
func NoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado"})
}