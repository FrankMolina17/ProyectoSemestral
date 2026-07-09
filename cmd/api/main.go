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
	// ── Configuración común ──
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sistema de Gestión y Control de Obras - API funcionando"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"estado":"ok"}`))
	})

	// ========================================
	// MÓDULO 1 - CATÁLOGO
	// ========================================
	catalogoStorage := storage.NuevoCatalogoStorage()
	catalogoStorage.Seed()

	materialSvc := services.NuevoMaterialService(catalogoStorage)
	manoObraSvc := services.NuevoManoObraService(catalogoStorage)
	equipoSvc := services.NuevoEquipoService(catalogoStorage)
	precioSvc := services.NuevoPreciosService(catalogoStorage)

	mh := handlers.NuevoMaterialHandler(materialSvc)
	mob := handlers.NuevoManoObraHandler(manoObraSvc)
	eh := handlers.NuevoEquipoHandler(equipoSvc)
	ph := handlers.NuevoPrecioHandler(precioSvc)

	authServiceCatalogo := services.NuevaAutenticacionService(catalogoStorage)

	// ========================================
	// MÓDULO 2 - PROFORMAS (tu módulo)
	// ========================================
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "proforma.db"
	}

	recursos, err := storage.InicializarModulo2(dsn)
	if err != nil {
		log.Fatalf("error inicializando módulo 2: %v", err)
	}
	defer recursos.Cerrar()

	proformaService := services.NuevoProformaService(recursos.ProformaRepo)
	authServiceProforma := services.NuevoAuthService(recursos.UsuarioStore)

	proformaHandler := handlers.NuevoHandler(proformaService)
	authHandler := handlers.NuevoAuthHandler(authServiceProforma)

	// ========================================
	// RUTAS
	// ========================================

	// Auth común
	r.Post("/api/v1/usuarios/registrar", authHandler.Registrar) // o serverC.RegistrarUser si prefieres
	r.Post("/api/v1/auth/login", authHandler.Login)

	// Rutas Catálogo (Módulo 1)
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

	// Rutas Proformas (Módulo 2)
	r.Group(func(r chi.Router) {
		r.Use(middleware.VerificarJWT(authServiceProforma))

		r.Route("/api/v1/proformas", func(r chi.Router) {
			r.Post("/", proformaHandler.CrearProforma)
			r.Get("/", proformaHandler.ObtenerTodos)
			r.Get("/{id}", proformaHandler.ObtenerPorID)
			r.Put("/{id}", proformaHandler.ActualizarProforma)
			r.Delete("/{id}", proformaHandler.EliminarProforma)
			r.Post("/{id}/items", proformaHandler.AgregarItem)
			r.Get("/{id}/items", proformaHandler.ObtenerItems)
			r.Put("/{id}/aprobar", proformaHandler.AprobarProforma)
			r.Get("/{id}/resumen", proformaHandler.ObtenerResumen)
			r.Post("/{id}/notas", proformaHandler.AgregarNota)
			r.Get("/{id}/notas", proformaHandler.ObtenerNotas)
		})

		r.Route("/api/v1/clientes", func(r chi.Router) {
			r.Post("/", proformaHandler.CrearCliente)
			r.Get("/", proformaHandler.ObtenerClientes)
			r.Get("/{id}", proformaHandler.ObtenerClientePorID)
			r.Put("/{id}", proformaHandler.ActualizarCliente)
			r.Delete("/{id}", proformaHandler.EliminarCliente)
		})
	})

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("API escuchando en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}