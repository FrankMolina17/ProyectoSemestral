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
	contextoClaveKey = contextoClave("usuarioID")
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
			if token == "" {
				NoAutorizado(w)
				return
			}
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

type claveContexto string

const ClaveusuarioID claveContexto = "userID"

// Middleware de autenticación basado en JWT.
func Autenticacion(s *services.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)

			if len(partes) != 2 || partes[0] != "Bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			userID, err := s.VerificarToken(partes[1])
			if err != nil {
				responderNoAutor(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaveusuarioID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado"})
}

// responderNoAutor responde con un 401 y un JSON de error.
func responderNoAutor(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":"Token inválido"}`))
}
