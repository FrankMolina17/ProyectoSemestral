package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"strings"
)

type IncidenciaService struct {
	repo storage.IncidenciaRepository
}

func NuevaIncidenciaService(repo storage.IncidenciaRepository) *IncidenciaService {
	return &IncidenciaService{repo: repo}
}

func (s *IncidenciaService) Listar() []models.Incidencia {
	return s.repo.ListarIncidencias()
}

func (s *IncidenciaService) Obtener(id int) (models.Incidencia, error) {
	c, ok := s.repo.BuscarIncidenciaPorID(id)
	if !ok {
		return models.Incidencia{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *IncidenciaService) ObtenerPorEntidad(id int, tipo string) (models.Incidencia, error) {
	c, ok := s.repo.BuscarIncidenciaPorEntidad(id, tipo)
	if !ok {
		return models.Incidencia{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *IncidenciaService) CrearIncidencia(c models.Incidencia) (models.Incidencia, error) {
	return s.repo.CrearIncidencia(c), nil
}

func (s *IncidenciaService) ActualizarIncidencia(id int, c models.Incidencia) (models.Incidencia, error) {
	if err := ValidacionIncidencia(c); err != nil {
		return models.Incidencia{}, err
	}

	actualizado, ok := s.repo.ActualizarIncidencia(id, c)
	if !ok {
		return models.Incidencia{}, ErrNoEncontrado
	}
	return actualizado, nil
}

func (s *IncidenciaService) BorrarIncidencia(id int) error {
	if !s.repo.BorrarIncidencia(id) {
		return ErrNoEncontrado
	}
	return nil
}

func ValidacionIncidencia(c models.Incidencia) error {
	if strings.TrimSpace(c.Titulo) == "" {
		return ErrTituloIncidenciaVacio
	}
	if strings.TrimSpace(c.Descripcion) == "" {
		return ErrDescripcionIncidenciaVacio
	}
	if strings.TrimSpace(c.Estado) == "" {
		return ErrEstadoIncidenciaVacio
	}
	return nil
}
