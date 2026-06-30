package storage

import (
	"time"

	"gorm.io/gorm"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

type UsuarioGORM struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioGORM {
	return &UsuarioGORM{db: db}
}

func (r *UsuarioGORM) CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error) {
	u := models.Usuario{
		Email:        in.Email,
		PasswordHash: in.Password,
		CreatedAt:    time.Now(),
	}
	if err := r.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}