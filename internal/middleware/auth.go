package middleware

import (
	"net/http"
	"strings"

	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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

			_, err := authSvc.VerificarToken(partes[1])
			if err != nil {
				http.Error(w, `{"error":"token inválido o expirado"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
