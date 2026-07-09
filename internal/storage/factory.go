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

type Recursos struct {
	Almacen      *Storage
	Usuarios     *Storage
	BackendUsado string
	Cerrar       func() error
}

func Inicializar(rutaDB string) (*Recursos, error) {
	_ = rutaDB
	almacen := New()
	almacen.Seed()
	cerrar := func() error { return nil }
	return &Recursos{
		Almacen:      almacen,
		Usuarios:     almacen,
		BackendUsado: "memoria",
		Cerrar:       cerrar,
	}, nil
}

type RecursosModulo2 struct {
	ProformaRepo repository.ProformaRepository
	UsuarioStore *UsuarioStorage
	Cerrar       func() error
	BackendUsado string
}

func InicializarModulo2(dsn string) (*RecursosModulo2, error) {
	driver := os.Getenv("DB_DRIVER")
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

	if err := db.AutoMigrate(
		&models.Proforma{},
		&models.ProformaItem{},
		&models.Cliente{},
		&models.NotaProforma{},
	); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	proformaRepo := repository.NuevoProformaRepository(db)
	usuarioStore := NuevoUsuarioStorage()

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
