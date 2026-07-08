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
	cfg := config.LoadConfig()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	recursos, err := storage.Inicializar(cfg.RutaDB, cfg.Backend, cfg.Driver, cfg.DSN)
	if err != nil {
		return err
	}
	log.Printf("Backend de almacenamiento activo: %s (ruta: %s)", recursos.BackendUsado, cfg.RutaDB)
	defer func() {
		if err := recursos.Cerrar(); err != nil {
			log.Printf("Backend cerrado (%s): %v", recursos.BackendUsado, err)
		}
	}()

	materialsvc := services.NewMaterialService(recursos.Almacen)
	manoobraSvc := services.NewManoObraService(recursos.ManoObra)
	equiposvc := services.NewEquipoService(recursos.Equipos)
	preciosSvc := services.NewPreciosService(recursos.Precios)
	authSvc := services.NuevaAutenticacionService(recursos.Usuarios, services.AuthOptions{
		Secreto:  []byte(cfg.JWTSecreto),
		Duracion: cfg.JWTDuracion,
	})

	serverC := handlers.NewServerC(manoobraSvc, materialsvc, equiposvc, preciosSvc, authSvc)

	mh := handlers.NewMaterialHandler(materialsvc)
	mob := handlers.NewManoObraHandler(manoobraSvc)
	eh := handlers.NewEquipoHandler(equiposvc)
	ph := handlers.NewPrecioHandler(preciosSvc)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Sistem-Inte-Gestion-Control-Obras/internal/routes - Servidor funcionando correctamente"))
	})

	r.Post("/api/v1/usuarios/registrar", serverC.RegistrarUser)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", serverC.LoginUser)
	})
	r.Post("/api/v1/usuarios/login", serverC.LoginUser)

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

	// 6. Servidor HTTP configurado por Options (puerto + timeouts desde config).
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	// 7. Contexto que se cancela al recibir Ctrl+C o SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 8. Arrancar el servidor en una goroutine para no bloquear la espera de la senal.
	errServidor := make(chan error, 1)
	go func() {
		log.Printf("Servidor escuchando en http://localhost%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	// 9. Esperar: o el servidor falla, o llega la senal de apagado.
	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Senal de apagado recibida, cerrando ordenadamente...")
	}

	// 10. Graceful shutdown: hasta 10s para terminar las requests en curso.
	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido limpiamente.")
	return nil
}
