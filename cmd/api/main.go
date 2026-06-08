package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

    "Sistem-Inte-Gestion-Control-Obras/internal/handlers"
    "Sistem-Inte-Gestion-Control-Obras/internal/middleware"
    "Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// modulo 1: Catalogo de Recursos
func main() {

	// Inicializar store y seed de datos
	s := storage.New()
	s.Seed()

	// Instanciar handlers con el store compartido
	mh  := handlers.NewMaterialHandler(s)
	mob := handlers.NewManoObraHandler(s)
	eh  := handlers.NewEquipoHandler(s)
	ph  := handlers.NewPrecioHandler(s)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1/catalogo", func(r chi.Router) {

		// JWT requerido en todos los endpoints del catálogo
		r.Use(middleware.AuthJWT)

		// material (/api/v1/catalogo/material)
		r.Get("/material", mh.Lista)
		r.Get("/material/{id}", mh.ObtenerUnMaterial)
		r.Post("/material", mh.CrearUnMaterial)    // valida nombre, unidad y precio>0, unicidad nombre+unidad
		r.Put("/material/{id}", mh.ActulizarUnMaterial)
		r.Delete("/material/{id}", mh.BorrarUnMaterial)

		// mano de obra (/api/v1/catalogo/manoobra)
		r.Get("/manoobra", mob.List)
		r.Get("/manoobra/{id}", mob.GetByID)
		r.Post("/manoobra", mob.Create)   // valida descripción, categoría en oficial/ayudante/especialista, unidad en hora/día/jornal, costo > 0
		r.Put("/manoobra/{id}", mob.Replace)
		r.Delete("/manoobra/{id}", mob.Delete)

		// equipo (/api/v1/catalogo/equipo)
		r.Get("/equipo", eh.ListaUnEquipo)
		r.Get("/equipo/{id}", eh.ObtenerporUnIDEquipo)
		r.Post("/equipo", eh.CrearUnEquipo)
		r.Put("/equipo/{id}", eh.ActualizarUnEquipo)
		r.Patch("/equipo/{id}/disponibilidad", eh.PatchDisponibilidad)
		r.Delete("/equipo/{id}", eh.BorrarUnEquipo)

		// precios (/api/v1/catalogo/precio)
		r.Post("/precio", ph.Create)
		r.Get("/precio/{tipo}/{id}", ph.Historial)         // historial completo de un recurso
		r.Get("/precio/{tipo}/{id}/vigente", ph.Vigente)   // precio más reciente con fecha_vigencia ≤ hoy
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
