package main

import (
	"net/http"

	"log"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Ruta de prueba
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Sistem-Inte-Gestion-Control-Obras/internal/routes - Servidor funcionando correctamente"))
	})

	// Inicializar storage y cargar datos de prueba
	s := storage.New()
	s.Seed()

	//  handlers y conectar con el storage
	mh := handlers.NewMaterialHandler(s)
	mob := handlers.NewManoObraHandler(s)
	eh := handlers.NewEquipoHandler(s)
	ph := handlers.NewPrecioHandler(s)

	r.Route("/api/v1/catalogo", func(r chi.Router) {

		// material (/api/v1/catalogo/material)
		r.Get("/material", mh.ListandoMateriales) //esto se conecta con el storage
		r.Get("/material/{id}", mh.ObtenerMaterialPorID)
		r.Post("/material", mh.CreandoMaterial) // valida nombre, unidad y precio>0, unicidad nombre+unidad
		r.Put("/material/{id}", mh.ActulizarUnMaterial)
		r.Delete("/material/{id}", mh.BorrarUnMaterial)

		// mano de obra (/api/v1/catalogo/manoobra)
		r.Get("/manoobra", mob.ListarUnaManoObra)
		r.Get("/manoobra/{id}", mob.ObtenerUnaManoObraPorID)
		r.Post("/manoobra", mob.CreandoUnaManoObra) // valida descripción, categoría en oficial/ayudante/especialista, unidad en hora/día/jornal, costo > 0
		r.Put("/manoobra/{id}", mob.ActualizadoUnaManoObra)
		r.Delete("/manoobra/{id}", mob.BorrandoUnaManoObra)

		// equipo (/api/v1/catalogo/equipo)
		r.Get("/equipo", eh.ListandoLosEquipos)
		r.Get("/equipo/{id}", eh.ObtenerUnEquipoPorID)
		r.Post("/equipo", eh.CrearUnEquipo)
		r.Put("/equipo/{id}", eh.ActualizarUnEquipo)
		r.Delete("/equipo/{id}", eh.BorrarUnEquipo)

		// precios (/api/v1/catalogo/precio)
		r.Get("/precio", ph.ListarDeLosPrecios)
		r.Post("/precio", ph.CrearUnPrecio)
		r.Get("/precio/{tipo}/{recursoID}/vigente", ph.PrecioVigentePorRecurso)
		r.Get("/precio/{tipo}/{recursoID}", ph.HistorialPorRecurso)
		r.Get("/precio/{id}", ph.ObtenerUnPrecioPorID)
		r.Put("/precio/{id}", ph.ActualizarUnPrecio)
		r.Delete("/precio/{id}", ph.BorrarUnPrecio)
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

	const addr = ":8080"
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
