// Package router registra todas las rutas del Módulo 1 · Catálogo de Recursos.
package repository

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

// New construye y devuelve el Chi router raíz listo para escuchar.
func New(s *storage.Storage) chi.Router {
	r := chi.NewRouter()

	// ── Middlewares globales ──────────────────────────────────────────────────
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// ── Subrouter /api/v1/catalogo ────────────────────────────────────────────
	r.Route("/api/v1/catalogo", func(r chi.Router) {

		// Todos los endpoints del catálogo requieren JWT
		r.Use(middleware.AuthJWT)

		// Instanciar handlers
		mh := handlers.NewMaterialHandler(s)
		mob := handlers.NewManoObraHandler(s)
		eh := handlers.NewEquipoHandler(s)
		ph := handlers.NewPrecioHandler(s)

		// ── /materiales ───────────────────────────────────────────────────────
		r.Route("/materiales", func(r chi.Router) {
			r.Get("/", mh.Lista)      // GET la lista de materiales
			r.Post("/", mh.CrearUnMaterial)   // POST /crear material
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", mh.ObtenerUnMaterial)    // GET    /materiales/:id
				r.Put("/", mh.ActulizarUnMaterial)    // PUT    /materiales/:id
				r.Delete("/", mh.BorrarUnMaterial)  // DELETE /materiales/:id
			})
		})

		// ── /mano-obra ────────────────────────────────────────────────────────
		r.Route("/mano-obra", func(r chi.Router) {
			r.Get("/", mob.List)     // GET  /mano-obra?categoria=
			r.Post("/", mob.Create)  // POST /mano-obra

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", mob.GetByID)    // GET    /mano-obra/:id
				r.Put("/", mob.Replace)    // PUT    /mano-obra/:id
				r.Delete("/", mob.Delete)  // DELETE /mano-obra/:id
			})
		})

		// ── /equipos ──────────────────────────────────────────────────────────
		r.Route("/equipos", func(r chi.Router) {
			r.Get("/", eh.ListaUnEquipo)    // GET  /equipos?disponible=true&tipo=
			r.Post("/", eh.CrearUnEquipo) // POST /equipos

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", eh.ObtenerporUnIDEquipo)                           // GET    /equipos/:id
				r.Put("/", eh.ActualizarUnEquipo)                           // PUT    /equipos/:id
				r.Delete("/", eh.BorrarUnEquipo)                         // DELETE /equipos/:id
				r.Patch("/disponibilidad", eh.PatchDisponibilidad) // PATCH  /equipos/:id/disponibilidad
			})
		})

		// ── /precios ──────────────────────────────────────────────────────────
		r.Route("/precios", func(r chi.Router) {
			r.Post("/", ph.Create) // POST /precios

			// GET /precios/{tipo}/{id}          → historial completo
			// GET /precios/{tipo}/{id}/vigente  → último precio vigente
			r.Route("/{tipo}/{id}", func(r chi.Router) {
				r.Get("/", ph.Historial)
				r.Get("/vigente", ph.Vigente)
			})
		})
	})

	return r
}
