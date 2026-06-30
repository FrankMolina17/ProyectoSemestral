package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/shopspring/decimal"
)

type PreciosService struct {
	repo storage.PrecioRecursoRepository
}

func NewPreciosService(repo storage.PrecioRecursoRepository) *PreciosService {
	return &PreciosService{
		repo: repo,
	}
}
//ListarPrecios() []*models.PrecioRecurso 
	
func (s *PreciosService) ListarPr() []*models.PrecioRecurso {
	return s.repo.ListarPrecios()
}

//ObtenerPrecio(id int) (*models.PrecioRecurso, error)
	
func (s *PreciosService) ObtenerPr(id int) (*models.PrecioRecurso, error) {
	p, err := s.repo.ObtenerPrecio(id)
	if err != nil {
		return &models.PrecioRecurso{}, ErrNoEncontrado
	}
	return p, nil
}

//CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error)

func (s *PreciosService) CrearPr(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	err := ValidarPr(in)
	if err != nil {
		return &models.PrecioRecurso{}, err
	}
	return s.repo.CrearPrecio(in)
}

//HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso
//es una función que se utiliza para obtener el historial de precios de un recurso.
func (s *PreciosService) HistorialPr(tipo string, recursoID int) []*models.PrecioRecurso {
	return s.repo.HistorialPrecios(tipo, recursoID)
}

//PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error)

//PrecioVigente es una función que se utiliza para obtener el precio vigente de un recurso.
func (s *PreciosService) PrecioVigentePr(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	return s.repo.PrecioVigente(tipo, recursoID)
}

//ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) 
func (s *PreciosService) ActualizarPr(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	if err := ValidarPr(in); err != nil {
		return &models.PrecioRecurso{}, err
	}
	actulizado, err := s.repo.ActualizarPrecio(id, in)
	if err != nil {
		return &models.PrecioRecurso{}, ErrNoEncontrado
	}
	return actulizado, nil
}

//existeRecurso(tipo string, id int) error 
//es una metodo que se utiliza para verificar si un recurso existe en la base de datos.
	
func (s *PreciosService) ExisteRecursoPr(tipo string, id int) error {
	return s.repo.ExisteRecurso(tipo, id)
}
//EliminarPrecio(id int) error 
//eliminar
	
func (s *PreciosService) EliminarPr(id int) error {
	if err := s.repo.EliminarPrecio(id); err != nil {
		return ErrNoEncontrado
	}
	return nil
}

//validar
	
func ValidarPr(in models.EntradaPrecioRecurso) error {
	if !models.RecursosTipos[in.RecursoTipo] {
		return ErrCategoriaNoPermitida
	}
	if in.RecursoID <= 0 {
		return ErrNoEncontrado
	}
	if in.Precio.LessThanOrEqual(decimal.Zero) {
		return ErrPrecioNegativo
	}
	if in.FechaVigencia.IsZero() {
		return ErrFechaVigenciaVacia
	}
	return nil
}


