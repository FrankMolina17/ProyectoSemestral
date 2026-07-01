package main

import (
	"net/http"

	"log"

	"Sistem-Inte-Gestion-Control-Obras/internal/handlers"
	"Sistem-Inte-Gestion-Control-Obras/internal/middleware"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

func main() {
	gdb, err := gorm.Open(sqlite.Open("incidencia.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}
	if err := gdb.AutoMigrate(&models.Incidencia{}, &models.Usuario{}); err != nil {
		log.Fatal("falló AutoMigrate: ", err)
	}
	almacenGorm := storage.NuevoAlmacenSQLite(gdb)

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(middleware.CORS)
	r.Use(chimw.Recoverer)

	// Ruta de prueba
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" Sistem-Inte-Gestion-Control-Obras/internal/routes - Servidor funcionando correctamente"))
	})

	var almacen storage.Almacen
	almacen = almacenGorm
	log.Println("Backend de almacenamiento: GORM")

	usuarioRepo := storage.NewUsuarioRepository(gdb)
	authService := services.NuevoAuthService(usuarioRepo)
	incidenciaService := services.NuevaIncidenciaService(almacen)
	servidor := handlers.NewServer(incidenciaService, authService)

	// Modulo 3 - Incidencias
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", servidor.Login)
		r.Post("/auth/register", servidor.Registrar)
	})

	r.Route("/api/v1/incidencias", func(r chi.Router) {
		r.Use(middleware.Autenticacion(*authService))
		r.Get("/", servidor.ObtenerIncidencias)
		r.Post("/", servidor.CrearIncidencia)
		r.Get("/{id}", servidor.ObtenerIncidenciaPorID)
		r.Get("/por/{tipo}/{id}", servidor.ObtenerIncidenciasPorEntidad)
		r.Put("/{id}", servidor.ActualizarIncidencia)
		r.Delete("/{id}", servidor.EliminarIncidencia)
	})

	const addr = ":8080"
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
