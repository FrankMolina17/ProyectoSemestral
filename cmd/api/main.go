package main

import (
	"net/http"

	"log"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Ruta de prueba
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Sistem-Inte-Gestion-Control-Obras/internal/routes - Servidor funcionando correctamente"))
	})

	proformaStore := storage.NuevoStorage()
	proformaHandler := handlers.NuevoHandler(proformaStore)

	// Módulo 2 — Proformas y Cálculo
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/proformas", proformaHandler.CrearProforma)
		r.Get("/proformas", proformaHandler.ObtenerTodos)
		r.Get("/proformas/{id}", proformaHandler.ObtenerPorID)
		r.Put("/proformas/{id}", proformaHandler.ActualizarProforma)
		r.Delete("/proformas/{id}", proformaHandler.EliminarProforma)

		// Items
		r.Post("/proformas/{id}/items", proformaHandler.AgregarItem)
		r.Get("/proformas/{id}/items", proformaHandler.ObtenerItems)
		r.Put("/proformas/{id}/aprobar", proformaHandler.AprobarProforma)
	})

	// Modulo 3 - Incidencias
	r.Route("/api/v1/incidencias", func(r chi.Router) {
		r.Post("/", handlers.CrearIncidenciaHandler)
		r.Get("/", handlers.ObtenerIncidenciasHandler)
		r.Get("/{id}", handlers.ObtenerIncidenciaPorIDHandler)
		r.Get("/por/{tipo}/{id}", handlers.ObtenerIncidenciasPorEntidadHandler)
		r.Put("/{id}", handlers.ActualizarIncidenciaHandler)
		r.Delete("/{id}", handlers.EliminarIncidenciaHandler)
	})

	const addr = ":8000"
	log.Printf("API escuchando en %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
