package storage

import (
	"fmt"
	"os"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// RecursosModulo2 agrupa todo lo que el módulo 2 necesita para funcionar.
type RecursosModulo2 struct {
	ProformaRepo  repository.ProformaRepository
	UsuarioStore  *UsuarioStorage
	Cerrar        func() error
	BackendUsado  string
}

// InicializarModulo2 centraliza el plumbing del módulo 2:
//  1. Abre GORM según el driver configurado (SQLite o PostgreSQL).
//  2. Ejecuta AutoMigrate para crear/actualizar las tablas.
//  3. Crea el repositorio de proformas.
//  4. Crea el storage de usuarios en memoria.
//  5. Expone una función Cerrar para el graceful shutdown.
func InicializarModulo2(dsn string) (*RecursosModulo2, error) {
	driver := os.Getenv("DB_DRIVER")

	// 1. Abrir GORM según el driver
	var db *gorm.DB
	var err error
	backendUsado := "sqlite"

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		backendUsado = "postgres"
	default:
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, fmt.Errorf("abrir base de datos: %w", err)
	}

	// 2. AutoMigrate — GORM es el dueño del esquema
	if err := db.AutoMigrate(
		&models.Proforma{},
		&models.ProformaItem{},
		&models.Cliente{},
		&models.NotaProforma{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	// 3. Repositorio de proformas
	proformaRepo := repository.NuevoProformaRepository(db)

	// 4. Storage de usuarios (siempre en memoria)
	usuarioStore := NuevoUsuarioStorage()

	// 5. Cierre ordenado
	cerrar := func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &RecursosModulo2{
		ProformaRepo: proformaRepo,
		UsuarioStore: usuarioStore,
		Cerrar:       cerrar,
		BackendUsado: backendUsado,
	}, nil
}
