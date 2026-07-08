package storage

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"time"

	"gorm.io/gorm"
)

type UsuarioGORM struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioGORM {
	return &UsuarioGORM{db: db}
}

func (r *UsuarioGORM) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	u.CreatedAt = time.Now()
	if err := r.db.Create(&u).Error; err != nil {
		return models.Usuario{}, err
	}
	return u, nil
}

func (r *UsuarioGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (r *UsuarioGORM) ListarUsuarios() []*models.Usuario {
	var usuarios []*models.Usuario
	r.db.Find(&usuarios)
	return usuarios
}

func (r *UsuarioGORM) ObtenerUsuarioPorID(id int) (*models.Usuario, bool) {
	var u models.Usuario
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, false
	}
	return &u, true
}
