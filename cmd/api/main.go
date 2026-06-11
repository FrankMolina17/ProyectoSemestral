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
    r.Post("/proformas", proformaHandler.CrearProforma)  // POST /api/v1/proformas — Crear proforma
    r.Get("/proformas", proformaHandler.ObtenerTodos)    // GET /api/v1/proformas — Listar proformas
    r.Get("/proformas/{id}", proformaHandler.ObtenerPorID) // GET /api/v1/proformas/{id} — Obtener proforma por ID
    r.Put("/proformas/{id}", proformaHandler.ActualizarProforma) // PUT /api/v1/proformas/{id} — Actualizar proforma
    r.Delete("/proformas/{id}", proformaHandler.EliminarProforma) // DELETE /api/v1/proformas/{id} — Eliminar proforma

	// Items
	r.Post("/proformas/{id}/items", proformaHandler.AgregarItem) // POST /api/v1/proformas/{id}/items — Agregar ítem
    r.Get("/proformas/{id}/items", proformaHandler.ObtenerItems) // GET /api/v1/proformas/{id}/items — Listar ítems
	r.Put("/proformas/{id}/aprobar", proformaHandler.AprobarProforma) // PUT /api/v1/proformas/{id}/aprobar — Aprobar proforma


})
	const addr = ":3000"
	log.Printf("API escuchando en %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}


