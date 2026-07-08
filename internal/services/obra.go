package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"strings"
)

type ObraService struct {
	repo storage.ObraRepository
}

func NuevaObraService(repo storage.ObraRepository) *ObraService {
	return &ObraService{repo: repo}
}

// Listar todas las obras
func (s *ObraService) Listar() []models.Obra {
	return s.repo.ListarObras()
}

// Obtener una obra por ID
func (s *ObraService) Obtener(id int) (models.Obra, error) {
	o, ok := s.repo.BuscarObraPorID(id)
	if !ok {
		return models.Obra{}, ErrNoEncontrado
	}
	return o, nil
}

// Crear obra con validaciones
func (s *ObraService) CrearObra(o models.Obra) (models.Obra, error) {
	if err := ValidacionObra(o); err != nil {
		return models.Obra{}, err
	}

	// Valor por defecto de estado
	if o.Estado == "" {
		o.Estado = "planificacion"
	}

	return s.repo.CrearObra(o), nil
}

// Actualizar obra
func (s *ObraService) ActualizarObra(id int, datos models.Obra) (models.Obra, error) {
	if err := ValidacionObra(datos); err != nil {
		return models.Obra{}, err
	}

	actualizado, ok := s.repo.ActualizarObra(id, datos)
	if !ok {
		return models.Obra{}, ErrNoEncontrado
	}
	return actualizado, nil
}

// Eliminar obra
func (s *ObraService) BorrarObra(id int) error {
	if !s.repo.BorrarObra(id) {
		return ErrNoEncontrado
	}
	return nil
}

// Validaciones
func ValidacionObra(o models.Obra) error {
	if strings.TrimSpace(o.Nombre) == "" {
		return ErrNombreObraVacio
	}
	if o.UserID <= 0 {
		return ErrUserIDRequerido
	}
	return nil
}
