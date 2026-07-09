package services

import "errors"

var (
	ErrEmailVacio            = errors.New("email y password son requeridos")
	ErrEmailEnUso            = errors.New("el email ya está registrado")
	ErrCredencialesInvalidos = errors.New("credenciales inválidas")
	ErrRecursoNoEncontrado   = errors.New("recurso no encontrado")
)
 