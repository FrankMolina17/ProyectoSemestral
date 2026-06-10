package main

import (
	"log"
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	proformaStore := storage.NuevoStorage()
	proformaHandler := handlers.NuevoHandler(proformaStore)

	r := chi.NewRouter()

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
	const addr = ":3000"
	log.Printf("API escuchando en %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}


