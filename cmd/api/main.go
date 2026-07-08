package main

import (
	"context"
	"errors"
	"log"
	"net/http"
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
	// 1. Inicializar almacenamiento

	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	recursos, err := storage.Inicializar(cfg.RutaDB)

	if err != nil {
		log.Fatalf("Error al inicializar almacenamiento: %v", err)
	}

	log.Printf("Backend de almacenamiento: %s", recursos.BackendUsado)

	// 2. Crear servicios (versión temporal por ahora)
	authService3 := services.NuevoAuthService(recursos.Usuarios)

	incidenciaService := services.NuevaIncidenciaService(recursos.Almacen)
	obraService := services.NuevaObraService(recursos.Almacen)

	// 3. Crear servidor
	servidor := handlers.NewServer(handlers.Deps{
		IncidenciaService: incidenciaService,
		ObraService:       obraService,
		Auth:              authService3,
	})

	// 4. Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

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

	// Modulo 3
	// Rutas
	r.Route("/api/v1", func(r chi.Router) {
		// Públicas
		r.Post("/auth/register3", servidor.Registrar)
		r.Post("/auth/login3", servidor.Login)

		// Protegidas
		r.Group(func(r chi.Router) {
			r.Use(middleware.Autenticacion(authService3))

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
