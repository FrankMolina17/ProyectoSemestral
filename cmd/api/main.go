package main

import (
	"log"
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/routes"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
)
func main() {
	proformaStore := storage.NewProformaStorage()
	proformaHandler := handlers.NewProformaHandler(proformaStore)

	r := chi.NewRouter()

	routes.RegisterRoutes(r)

	// Módulo 2 — Proformas y Cálculo
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/proformas", proformaHandler.CrearProforma)
		r.Get("/proformas", proformaHandler.ListarProformas)
		r.Get("/proformas/{id}", proformaHandler.ObtenerProforma)
		r.Put("/proformas/{id}", proformaHandler.ActualizarProforma)

		// Items
		r.Put("/proformas/{id}/aprobar", proformaHandler.CalcularProforma)
	})
	const addr = ":3000" //http://localhost:3000
	log.Printf("API escuchando en %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}


