package routes

import (
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, auth *services.AutenticacionService) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gestión de Obras e Incidencias - Funcionando"))
	})

	r.Route("/api/v1/incidencias", func(r chi.Router) {
		r.Post("/", handlers.CrearIncidenciaHandler)
		r.Get("/", handlers.ObtenerIncidenciasHandler)
		r.Get("/{id}", handlers.ObtenerIncidenciaPorIDHandler)
		r.Get("/por/{tipo}/{id}", handlers.ObtenerIncidenciasPorEntidadHandler)
		r.Put("/{id}", handlers.ActualizarIncidenciaHandler)
	})
}
