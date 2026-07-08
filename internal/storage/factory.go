package storage

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

// Recursos agrupa todo lo necesario
type Recursos struct {
	Almacen      Almacen
	Usuarios     UserRepository
	BackendUsado string
	Cerrar       func() error
}

// Inicializar crea los recursos de almacenamiento
func Inicializar(rutaDB string) (*Recursos, error) {
	gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("abrir GORM: %w", err)
	}

	if err := gdb.AutoMigrate(&models.Incidencia{}, &models.Usuario{}, &models.Obra{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}

	almacenGorm := NuevoAlmacenSQLite(gdb)

	usuarios := NewUsuarioRepository(gdb)

	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		Almacen:      almacenGorm,
		Usuarios:     usuarios,
		BackendUsado: "gorm",
		Cerrar:       cerrar,
	}, nil
}
