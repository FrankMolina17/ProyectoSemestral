package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// modulo 1: Catalogo de Recursos
func main() {

	// Inicializar storage y cargar datos de prueba
	s := storage.New()
	s.Seed()

	//  handlers y conectar con el storage
	mh  := handlers.NewMaterialHandler(s)
	mob := handlers.NewManoObraHandler(s)
	eh  := handlers.NewEquipoHandler(s)
	ph  := handlers.NewPrecioHandler(s)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api/v1/catalogo", func(r chi.Router) {

		// material (/api/v1/catalogo/material)
		r.Get("/material", mh.ListandoMateriales) //esto se conecta con el storage
		r.Get("/material/{id}", mh.ObtenerMaterialPorID)
		r.Post("/material", mh.CreandoMaterial)    // valida nombre, unidad y precio>0, unicidad nombre+unidad
		r.Put("/material/{id}", mh.ActulizarUnMaterial)
		r.Delete("/material/{id}", mh.BorrarUnMaterial)

		// mano de obra (/api/v1/catalogo/manoobra)
		r.Get("/manoobra", mob.ListarUnaManoObra)
		r.Get("/manoobra/{id}", mob.ObtenerUnaManoObraPorID)
		r.Post("/manoobra", mob.CreandoUnaManoObra)   // valida descripción, categoría en oficial/ayudante/especialista, unidad en hora/día/jornal, costo > 0
		r.Put("/manoobra/{id}", mob. ActualizadoUnaManoObra)
		r.Delete("/manoobra/{id}", mob.BorrandoUnaManoObra)

		// equipo (/api/v1/catalogo/equipo)
		r.Get("/equipo", eh.ListandoLosEquipos)
		r.Get("/equipo/{id}", eh.ObtenerUnEquipoPorID)
		r.Post("/equipo", eh.CrearUnEquipo)
		r.Put("/equipo/{id}", eh.ActualizarUnEquipo)
		r.Delete("/equipo/{id}", eh.BorrarUnEquipo)

		// precios (/api/v1/catalogo/precio)
		r.Get("/precio", ph.ListarDeLosPrecios)
		r.Post("/precio", ph.	CrearUnPrecio)
		r.Get("/precio/{id}", ph.ObtenerUnPrecioPorID)
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
