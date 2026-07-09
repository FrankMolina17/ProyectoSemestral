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
	// ── Módulo 1 — Catálogo (Melani) ──
	s := storage.NuevoCatalogoStorage()
	s.Seed()

	manoObraSvc := services.NuevoManoObraService(s)
	materialSvc := services.NuevoMaterialService(s)
	equipoSvc := services.NuevoEquipoService(s)
	precioSvc := services.NuevoPreciosService(s)

	mh := handlers.NuevoMaterialHandler(materialSvc)
	mob := handlers.NuevoManoObraHandler(manoObraSvc)
	eh := handlers.NuevoEquipoHandler(equipoSvc)
	ph := handlers.NuevoPrecioHandler(precioSvc)

	authServiceCatalogo := services.NuevaAutenticacionService(s)
	serverC := handlers.NuevoServerC(
		manoObraSvc,
		materialSvc,
		equipoSvc,
		precioSvc,
		authServiceCatalogo,
	)

	// ── Módulo 2 — Proformas (Franklin) — patrón Factory ──
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "proforma.db"
	}
	recursos, err := storage.InicializarModulo2(dsn)
	if err != nil {
		log.Fatalf("error inicializando módulo 2: %v", err)
	}
	defer recursos.Cerrar()

	log.Printf("Módulo 2 usando backend: %s", recursos.BackendUsado)

	proformaService := services.NuevoProformaService(recursos.ProformaRepo)
	authService := services.NuevoAuthService(recursos.UsuarioStore)

	proformaHandler := handlers.NuevoHandler(proformaService)
	authHandler := handlers.NuevoAuthHandler(authService)

	// ── Router ──
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sistema de Gestion y Control de Obras - API funcionando"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"estado":"ok","modulo":"proformas"}`))
	})

	// Auth Módulo 1
	r.Post("/api/v1/usuarios/registrar", serverC.RegistrarUser)
	r.Post("/aios/login", serverC.LoginUser)
	r.Route("/api/v1/auth/catalogo", func(r chi.Router) {
		r.Post("/login", serverC.LoginUser)
	})

	// Auth Módulo 2
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Registrar)
		r.Post("/login", authHandler.Login)
	})

	// Rutas Módulo 1
	r.Route("/api/v1/catalogo", func(r chi.Router) {
		r.Use(middleware.AuthJWT(authServiceCatalogo))

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

	// Rutas Módulo 2
	r.Group(func(r chi.Router) {
		r.Use(middleware.VerificarJWT(authService))

		r.Route("/api/v1", func(r chi.Router) {
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
	log.Printf("API escuchando en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
