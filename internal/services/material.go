package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/shopspring/decimal"
	"strings"
)

type MaterialService struct {
	repo storage.MaterialRepository
}

func NewMaterialService(repo storage.MaterialRepository) *MaterialService {
	return &MaterialService{
		repo: repo,
	}
}

func (s *MaterialService) Listado() []*models.Material {
	return s.repo.ListarMateriales()
}

func (s *MaterialService) ObtenerM(id int) (*models.Material, error) {
	p, ok := s.repo.ObtenerMateriales(id)
	if !ok {
		return nil, ErrNoEncontrado
	}
	return p, nil
}

func (s *MaterialService) CrearM(in models.EntradaMaterial) (*models.Material, error) {
	if err := ValidarM(in); err != nil {
		return nil, err
	}
	return s.repo.CrearMateriales(in)
}

func (s *MaterialService) ActualizarM(id int, in models.EntradaMaterial) (*models.Material, error) {
	if err := ValidarM(in); err != nil {
		return nil, err
	}
	actualizado, ok := s.repo.ActualizarMateriales(id, in)
	if !ok {
		return nil, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *MaterialService) EliminarM(id int) error {
	if !s.repo.EliminarMateriales(id) {
		return ErrNoEncontrado
	}
	return nil
}

func ValidarM(in models.EntradaMaterial) error {
	if strings.TrimSpace(in.Nombre) == "" {
		return ErrNombreVacio
	}
	if strings.TrimSpace(in.Unidad) == "" {
		return ErrUnidadVacia
	}
	if strings.TrimSpace(in.Descripcion) == "" {
		return ErrDescripcionVacia
	}
	if !models.UnidadesPermitidas[in.Unidad] {
		return ErrUnidadNoPermitida
	}
	if in.PrecioReferencia.LessThanOrEqual(decimal.Zero) {
		return ErrPrecioReferencialInvalido
	}
	return nil
}
