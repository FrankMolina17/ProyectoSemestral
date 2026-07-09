package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/config"
	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/httpserver"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	// ── Módulo 1 — Catálogo ──
	recursos, err := storage.Inicializar(cfg.RutaDB)
	if err != nil {
		return err
	}
	defer recursos.Cerrar()

	materialSvc := services.NewMaterialService(recursos.Almacen)
	manoObraSvc := services.NewManoObraService(recursos.Almacen)
	equipoSvc := services.NewEquipoService(recursos.Almacen)
	precioSvc := services.NewPreciosService(recursos.Almacen)
	authSvc := services.NuevaAutenticacionService(recursos.Usuarios, services.AuthOptions{
		Secreto:  cfg.JWTSecreto,
		Duracion: cfg.JWTDuracion,
	})

	serverC := handlers.NewServerC(manoObraSvc, materialSvc, equipoSvc, precioSvc, authSvc)

	mh := handlers.NewMaterialHandler(materialSvc)
	mob := handlers.NewManoObraHandler(manoObraSvc)
	eh := handlers.NewEquipoHandler(equipoSvc)
	ph := handlers.NewPrecioHandler(precioSvc)

	// ── Módulo 2 — Proformas ──
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "proforma.db"
	}

	recursos2, err := storage.InicializarModulo2(dsn)
	if err != nil {
		return err
	}
	defer recursos2.Cerrar()

	proformaService := services.NuevoProformaService(recursos2.ProformaRepo)
	authServiceProforma := services.NuevoAuthService(recursos2.UsuarioStore)

	proformaHandler := handlers.NuevoHandler(proformaService)
	authHandler := handlers.NuevoAuthHandler(authServiceProforma)

	// ── Router ──
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gesti\u00f3n de Obras e Incidencias - Funcionando"))
	})

	// Auth Catálogo
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", serverC.LoginUser)
	})

	// Auth Proformas
	r.Post("/api/v1/auth/register", authHandler.Registrar)
	r.Post("/api/v1/auth/login", authHandler.Login)

	// ── Rutas Catálogo ──
	r.Route("/api/v1/catalogo", func(r chi.Router) {
		r.Use(middleware.AuthJWT(authSvc))

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

	// ── Rutas Proformas ──
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

	// ── Incidencias (desde routes.go) ──
	r.Route("/api/v1/incidencias", func(r chi.Router) {
		r.Post("/", handlers.CrearIncidenciaHandler)
		r.Get("/", handlers.ObtenerIncidenciasHandler)
		r.Get("/{id}", handlers.ObtenerIncidenciaPorIDHandler)
		r.Get("/por/{tipo}/{id}", handlers.ObtenerIncidenciasPorEntidadHandler)
		r.Put("/{id}", handlers.ActualizarIncidenciaHandler)
	})

	// ── Servidor HTTP ──
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Se\u00f1al de apagado recibida, cerrando ordenadamente...")
	}

	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido limpiamente.")
	return nil
}
