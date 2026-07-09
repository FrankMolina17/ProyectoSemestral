package services

import (

	"github.com/shopspring/decimal"
	"strings"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	
)

type EquipoService struct {
	repo storage.EquipoRepository
}

func NewEquipoService(repo storage.EquipoRepository) *EquipoService {
	return &EquipoService{
		repo: repo,
	}
}

//listar
func (s *EquipoService) ListadoE() []*models.Equipo {
	return s.repo.ListarEquipos()
}

//obtener
func (s *EquipoService) ObtenerE(id int) (*models.Equipo, error) {
	e, err := s.repo.ObtenerEquipo(id)
	if err != nil {
		return &models.Equipo{}, ErrNoEncontrado
	}
	return e, nil
}

//crear
func (s *EquipoService) CrearE(in models.EntradaEquipo) (*models.Equipo, error) {
	if err := ValidarE(in); err != nil {
		return &models.Equipo{}, err
	}
	return s.repo.CrearEquipo(in)
}
//Actualizar
func (s *EquipoService) ActualizarE(id int, in models.EntradaEquipo) (*models.Equipo, error) {
	if err := ValidarE(in); err != nil {
		return &models.Equipo{}, err
	}
	actualizado, err := s.repo.ActualizarEquipo(id, in)
	if err != nil {
		return &models.Equipo{}, ErrNoEncontrado
	}
	return actualizado, nil

}

func (s *EquipoService) EliminarE(id int) error {
	if err := s.repo.EliminarEquipo(id); err != nil {
		return ErrNoEncontrado
	}
	return nil
}

//validar
func ValidarE(in models.EntradaEquipo) error {

	if strings.TrimSpace(in.Nombre) == "" {
		return ErrNombreVacio
	}
	if !models.TiposEquipo[in.Tipo] {
		return ErrTipoVacio
	}
	if strings.TrimSpace(in.Unidad) == "" {
		return ErrUnidadVacia
	}
	if in.CostoHora.LessThanOrEqual(decimal.Zero) {
		return ErrPrecioReferencialInvalido
	}
	return nil
}

