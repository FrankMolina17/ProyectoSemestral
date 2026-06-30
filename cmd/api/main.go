package main

import (
	"net/http"

	"log"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Sistem-Inte-Gestion-Control-Obras/internal/routes - Servidor funcionando correctamente"))
	})

	s := storage.New()
	s.Seed()

	mh := handlers.NewMaterialHandler(s)
	mob := handlers.NewManoObraHandler(s)
	eh := handlers.NewEquipoHandler(s)
	ph := handlers.NewPrecioHandler(s)

	authService := services.NuevaAutenticacionService(s)
	serverC := handlers.NewServerC(
		services.NewManoObraService(s),
		services.NewMaterialService(s),
		services.NewEquipoService(s),
		services.NewPreciosService(s),
		authService,
	)

	r.Post("/api/v1/usuarios/registrar", serverC.RegistrarUser)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", serverC.LoginUser)
	})
	r.Post("/api/v1/usuarios/login", serverC.LoginUser)

	r.Route("/api/v1/catalogo", func(r chi.Router) {
		r.Use(middleware.AuthJWT(authService))

		r.Get("/material", mh.ListandoMateriales)
		r.Get("/material/{id}", mh.ObtenerMaterialPorID)
		r.Post("/material", mh.CreandoMaterial)
		r.Put("/material/{id}", mh.ActulizarUnMaterial)
		r.Delete("/material/{id}", mh.BorrarUnMaterial)

		r.Get("/manoobra", mob.ListarUnaManoObra)
		r.Get("/manoobra/{id}", mob.ObtenerUnaManoObraPorID)
		r.Post("/manoobra", mob.CreandoUnaManoObra)
		r.Put("/manoobra/{id}", mob.ActualizadoUnaManoObra)
		r.Delete("/manoobra/{id}", mob.BorrandoUnaManoObra)

		r.Get("/equipo", eh.ListandoLosEquipos)
		r.Get("/equipo/{id}", eh.ObtenerUnEquipoPorID)
		r.Post("/equipo", eh.CrearUnEquipo)
		r.Put("/equipo/{id}", eh.ActualizarUnEquipo)
		r.Delete("/equipo/{id}", eh.BorrarUnEquipo)

		r.Get("/precio", ph.ListarDeLosPrecios)
		r.Post("/precio", ph.CrearUnPrecio)
		r.Get("/precio/{tipo}/{recursoID}/vigente", ph.PrecioVigentePorRecurso)
		r.Get("/precio/{tipo}/{recursoID}", ph.HistorialPorRecurso)
		r.Get("/precio/{id}", ph.ObtenerUnPrecioPorID)
		r.Put("/precio/{id}", ph.ActualizarUnPrecio)
		r.Delete("/precio/{id}", ph.BorrarUnPrecio)
	})

	const addr = ":8080"
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
