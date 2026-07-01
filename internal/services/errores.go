package services

import "errors"

var (
	ErrNombreVacio                = errors.New("Nombre no puede estar vacío")
	ErrPrecioNegativo             = errors.New("Valor negativo inválido")
	ErrNoEncontrado               = errors.New("Recurso no encontrado")
	ErrEmailEnUso                 = errors.New("Email ya en uso")
	ErrCredencialesInvalidos      = errors.New("Email o contraseña incorrecta")
	ErrTituloIncidenciaVacio      = errors.New("El titulo de la incidencia no puede estar vacio")
	ErrDescripcionIncidenciaVacio = errors.New("La descripcion de la incidencia no puede estar vacia")
	ErrEstadoIncidenciaVacio      = errors.New("El estado de la incidencia no puede estar vacio")
)
