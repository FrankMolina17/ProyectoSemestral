package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"Sistem-Inte-Gestion-Control-Obras/internal/config"
	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/httpserver"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	// 1. Inicializar almacenamiento
	recursos, err := storage.Inicializar(cfg.RutaDB)
	if err != nil {
		return err
	}
	defer func() { _ = recursos.Cerrar() }()

	log.Printf("Backend de almacenamiento: %s", recursos.BackendUsado)

	// 2. Crear servicios (versión temporal por ahora)
	authService := services.NuevoAuthService(recursos.Usuarios)

	incidenciaService := services.NuevaIncidenciaService(recursos.Almacen)
	obraService := services.NuevaObraService(recursos.Almacen)

	// 3. Crear servidor
	servidor := handlers.NewServer(handlers.Deps{
		IncidenciaService: incidenciaService,
		ObraService:       obraService,
		Auth:              authService,
	})

	// 4. Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// Rutas
	r.Route("/api/v1", func(r chi.Router) {
		// Públicas
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)

		// Protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.Autenticacion(authService))

			r.Route("/incidencias", func(r chi.Router) {
				r.Get("/", servidor.ObtenerIncidencias)
				r.Post("/", servidor.CrearIncidencia)
				r.Get("/{id}", servidor.ObtenerIncidenciaPorID)
				r.Get("/por/{tipo}/{id}", servidor.ObtenerIncidenciasPorEntidad)
				r.Put("/{id}", servidor.ActualizarIncidencia)
				r.Delete("/{id}", servidor.EliminarIncidencia)
			})

			r.Route("/obras", func(r chi.Router) {
				r.Get("/", servidor.ObtenerObras)
				r.Post("/", servidor.CrearObra)
				r.Get("/{id}", servidor.ObtenerObraPorID)
				r.Put("/{id}", servidor.ActualizarObra)
				r.Delete("/{id}", servidor.EliminarObra)
			})
		})
	})

	// 5. Servidor HTTP
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 6. Graceful Shutdown
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
		log.Println("Señal de apagado recibida, cerrando...")
	}

	ctxApagado, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}

	log.Println("Servidor detenido limpiamente.")
	return nil
}
