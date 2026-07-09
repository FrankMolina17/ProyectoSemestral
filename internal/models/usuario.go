package models

import (
	"errors"
	"time"
)

type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Rol          string    `json:"rol" gorm:"not null;default:cliente"`
	CreatedAt    time.Time `json:"creado_en" gorm:"autoCreateTime"`
}

type EntradaUsuario struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Rol      string `json:"rol,omitempty"`
}

func (u EntradaUsuario) ValidarUsuario() error {
	if u.Email == "" {
		return errors.New("campo 'email' requerido")
	}
	if len(u.Password) < 6 {
		return errors.New("password debe tener al menos 6 caracteres")
	}
	if u.Rol != "" && u.Rol != "admin" && u.Rol != "cliente" {
		return errors.New("rol inválido: debe ser 'admin' o 'cliente'")
	}
	return nil
}
