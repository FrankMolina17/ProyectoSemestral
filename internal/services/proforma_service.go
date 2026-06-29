package services

import (
	"errors"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"
)

var (
	ErrNombreRequerido     = errors.New("el campo nombre es requerido")
	ErrObraIDRequerido     = errors.New("el campo obra_id es requerido")
	ErrProformaYaAprobada  = errors.New("la proforma ya está aprobada")
	ErrDescripcionRequerida = errors.New("el campo descripcion es requerido")
	ErrCantidadInvalida    = errors.New("la cantidad debe ser mayor a 0")
	ErrPrecioInvalido      = errors.New("el precio debe ser mayor a 0")
	ErrRucRequerido        = errors.New("el campo ruc es requerido")
	ErrContenidoRequerido  = errors.New("el campo contenido es requerido")
)

type ProformaService struct {
	repo repository.ProformaRepository
}

func NuevoProformaService(repo repository.ProformaRepository) *ProformaService {
	return &ProformaService{repo: repo}
}

func (s *ProformaService) CrearProforma(p models.Proforma) (models.Proforma, error) {
	if p.Nombre == "" {
		return models.Proforma{}, ErrNombreRequerido
	}
	if p.ObraID == 0 {
		return models.Proforma{}, ErrObraIDRequerido
	}
	return s.repo.CrearProforma(p)
}

func (s *ProformaService) ObtenerPorID(id int) (models.Proforma, error) {
	return s.repo.ObtenerPorID(id)
}

func (s *ProformaService) ObtenerTodos() ([]models.Proforma, error) {
	return s.repo.ObtenerTodos()
}

func (s *ProformaService) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
	if datos.Nombre == "" {
		return models.Proforma{}, ErrNombreRequerido
	}
	return s.repo.ActualizarProforma(id, datos)
}

func (s *ProformaService) EliminarProforma(id int) error {
	return s.repo.EliminarProforma(id)
}

func (s *ProformaService) AprobarProforma(id int) (models.Proforma, error) {
	p, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return models.Proforma{}, err
	}
	if p.Estado == "aprobada" {
		return models.Proforma{}, ErrProformaYaAprobada
	}
	return s.repo.AprobarProforma(id)
}

func (s *ProformaService) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
	if item.Descripcion == "" {
		return models.ProformaItem{}, ErrDescripcionRequerida
	}
	if item.Cantidad <= 0 {
		return models.ProformaItem{}, ErrCantidadInvalida
	}
	if item.PrecioPromedio <= 0 {
		return models.ProformaItem{}, ErrPrecioInvalido
	}
	return s.repo.AgregarItem(item)
}

func (s *ProformaService) ObtenerItems(proformaID int) ([]models.ProformaItem, error) {
	return s.repo.ObtenerItems(proformaID)
}

func (s *ProformaService) AgregarNota(nota models.NotaProforma) (models.NotaProforma, error) {
	if nota.Contenido == "" {
		return models.NotaProforma{}, ErrContenidoRequerido
	}
	return s.repo.AgregarNota(nota)
}

func (s *ProformaService) ObtenerNotas(proformaID int) ([]models.NotaProforma, error) {
	return s.repo.ObtenerNotas(proformaID)
}

func (s *ProformaService) CrearCliente(c models.Cliente) (models.Cliente, error) {
	if c.Nombre == "" {
		return models.Cliente{}, ErrNombreRequerido
	}
	if c.Ruc == "" {
		return models.Cliente{}, ErrRucRequerido
	}
	return s.repo.CrearCliente(c)
}

func (s *ProformaService) ObtenerClientes() ([]models.Cliente, error) {
	return s.repo.ObtenerClientes()
}

func (s *ProformaService) ObtenerClientePorID(id int) (models.Cliente, error) {
	return s.repo.ObtenerClientePorID(id)
}

func (s *ProformaService) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error) {
	if datos.Nombre == "" {
		return models.Cliente{}, ErrNombreRequerido
	}
	return s.repo.ActualizarCliente(id, datos)
}

func (s *ProformaService) EliminarCliente(id int) error {
	return s.repo.EliminarCliente(id)
}

func (s *ProformaService) ObtenerResumen(id int) (map[string]interface{}, error) {
	p, err := s.repo.ObtenerPorID(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"proforma_id":    p.ID,
		"nombre":         p.Nombre,
		"estado":         p.Estado,
		"subtotal":       p.Subtotal,
		"ganancia":       p.Subtotal * p.PctGanancia,
		"imprevisto":     p.Subtotal * p.PctImprevisto,
		"total":          p.Total,
		"pct_ganancia":   p.PctGanancia,
		"pct_imprevisto": p.PctImprevisto,
	}, nil
}
