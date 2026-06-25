package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// ── Storage ──
	usuarioStore := storage.NuevoUsuarioStorage()
	proformaStore := storage.NuevoStorage()

	// ── Services ──
	authService := services.NuevoAuthService(usuarioStore)

	// ── Handlers ──
	proformaHandler := handlers.NuevoHandler(proformaStore)
	authHandler := handlers.NuevoAuthHandler(authService)

	// ── Router ──
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"estado":"ok","modulo":"proformas"}`))
	})

	// ── Rutas públicas ──
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Registrar)
		r.Post("/login", authHandler.Login)
	})

	// ── Rutas protegidas — requieren JWT ──
	r.Group(func(r chi.Router) {
		r.Use(middleware.VerificarJWT(authService))

		r.Route("/api/v1", func(r chi.Router) {
			// Proformas
			r.Post("/proformas", proformaHandler.CrearProforma)
			r.Get("/proformas", proformaHandler.ObtenerTodos)
			r.Get("/proformas/{id}", proformaHandler.ObtenerPorID)
			r.Put("/proformas/{id}", proformaHandler.ActualizarProforma)
			r.Delete("/proformas/{id}", proformaHandler.EliminarProforma)
			r.Post("/proformas/{id}/items", proformaHandler.AgregarItem)
			r.Get("/proformas/{id}/items", proformaHandler.ObtenerItems)
			r.Put("/proformas/{id}/aprobar", proformaHandler.AprobarProforma)
			r.Get("/proformas/{id}/resumen", proformaHandler.ObtenerResumen)
			r.Post("/proformas/{id}/notas", proformaHandler.AgregarNota)
			r.Get("/proformas/{id}/notas", proformaHandler.ObtenerNotas)

			// Clientes
			r.Post("/clientes", proformaHandler.CrearCliente)
			r.Get("/clientes", proformaHandler.ObtenerClientes)
			r.Get("/clientes/{id}", proformaHandler.ObtenerClientePorID)
			r.Put("/clientes/{id}", proformaHandler.ActualizarCliente)
			r.Delete("/clientes/{id}", proformaHandler.EliminarCliente)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("M2 Proformas escuchando en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}