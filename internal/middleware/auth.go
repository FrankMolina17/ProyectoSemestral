package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

type contextoClave string

const (
	contextoClaveKey    = contextoClave("usuarioID")
	contextoClaveRol    = contextoClave("rol")
)

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
			ctx = context.WithValue(ctx, contextoClaveRol, claims.Rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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

			claims, err := authSvc.VerificarToken(partes[1])
			if err != nil {
				http.Error(w, `{"error":"token inválido o expirado"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextoClaveKey, claims.UsuarioID)
			ctx = context.WithValue(ctx, contextoClaveRol, claims.Rol)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequerirRol(rol string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRol, ok := r.Context().Value(contextoClaveRol).(string)
			if !ok || userRol != rol {
				Prohibido(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado"})
}

func Prohibido(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "No tienes permisos para acceder a este recurso"})
}
