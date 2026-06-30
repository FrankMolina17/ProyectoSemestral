package models

import (
	"errors"
	"time"
)

type Usuario struct {
	ID           int       `json:"id"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

type EntradaUsuario struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u EntradaUsuario) ValidarUsuario() error {
	if u.Email == "" {
		return errors.New("campo 'email' requerido")
	}
	if len(u.Password) < 6 {
		return errors.New("password debe tener al menos 6 caracteres")
	}
	return nil
}
