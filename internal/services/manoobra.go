package services

import (

	"github.com/shopspring/decimal"
	"strings"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	
)

type ManoObraServise struct{
	repo storage.ManoObraRepository
}

func NewManoObraService(repo storage.ManoObraRepository) *ManoObraServise {
	return &ManoObraServise{
		repo: repo,
	}
}
//listar
func (s *ManoObraServise) ListadoMa() []*models.ManoObra {
	return s.repo.ListarManoObra()
}

//obtener
func (s *ManoObraServise) ObtenerMa(id int) (*models.ManoObra, error) {
	m, ok := s.repo.ObtenerManoObra(id)
	if !ok {
		return &models.ManoObra{}, ErrNoEncontrado
	}
	return m, nil
}

//crear
func (s *ManoObraServise) CrearMa(in models.EntradaManoObra) (*models.ManoObra, error) {
	if err := ValidarMa(in); err != nil {
		return &models.ManoObra{}, err
	}
	return s.repo.CrearManoObra(in)
}

//actualizar
func (s *ManoObraServise) ActualizarMa(id int, in models.EntradaManoObra) (*models.ManoObra, error) {
	if err := ValidarMa(in); err != nil {
		return &models.ManoObra{}, err
	}
	actualizado, ok := s.repo.ActualizarManoObra(id, in)
	if !ok {
		return &models.ManoObra{}, ErrNoEncontrado
	}
	return actualizado, nil
}

//eliminar
func (s *ManoObraServise) EliminarMa(id int) bool {
	return s.repo.EliminarManoObra(id)
}

//validar
func ValidarMa(in models.EntradaManoObra) error {
	if strings.TrimSpace(in.Descripcion) == "" {
		return ErrDescripcionVacia
	}
	if strings.TrimSpace(in.Categoria) == "" {
		return ErrCategoriaNoPermitida
	
	}
	if strings.TrimSpace(in.Unidad) == "" {
		return ErrUnidadVacia
	}
	if in.CostoReferencia.LessThanOrEqual(decimal.Zero) {
		return ErrPrecioReferencialInvalido
	}
	return nil
}

