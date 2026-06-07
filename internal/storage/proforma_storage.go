package storage

import (
	"errors"
	"sync"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

const tasaImpuesto = 0.18

var (
	ErrProformaNoEncontrada = errors.New("proforma no encontrada")
	ErrNombreRequerido      = errors.New("el nombre de la proforma es requerido")
)

// ProformaStorage almacena proformas en memoria.
type ProformaStorage struct {
	mu        sync.RWMutex
	proformas map[int]*models.Proforma
	nextID    int
	nextItem  int
}

// NewProformaStorage crea un almacenamiento vacío en memoria.
func NewProformaStorage() *ProformaStorage {
	return &ProformaStorage{
		proformas: make(map[int]*models.Proforma),
		nextID:    1,
		nextItem:  1,
	}
}

// Crear guarda una nueva proforma y calcula sus totales.
func (s *ProformaStorage) Crear(p *models.Proforma) (*models.Proforma, error) {
	if p.Nombre == "" {
		return nil, ErrNombreRequerido
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	copia := *p
	copia.ID = s.nextID
	s.nextID++

	for i := range copia.Items {
		copia.Items[i].ID = s.nextItem
		s.nextItem++
	}

	if copia.Estado == "" {
		copia.Estado = "borrador"
	}

	calcularTotales(&copia)
	s.proformas[copia.ID] = &copia

	return &copia, nil
}

// ObtenerPorID devuelve una proforma por su identificador.
func (s *ProformaStorage) ObtenerPorID(id int) (*models.Proforma, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.proformas[id]
	if !ok {
		return nil, ErrProformaNoEncontrada
	}

	copia := *p
	return &copia, nil
}

// Listar devuelve todas las proformas registradas.
func (s *ProformaStorage) Listar() []models.Proforma {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resultado := make([]models.Proforma, 0, len(s.proformas))
	for _, p := range s.proformas {
		resultado = append(resultado, *p)
	}

	return resultado
}

// Actualizar modifica una proforma existente y recalcula sus totales.
func (s *ProformaStorage) Actualizar(id int, actualizada models.Proforma) (*models.Proforma, error) {
	if actualizada.Nombre == "" {
		return nil, ErrNombreRequerido
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.proformas[id]; !ok {
		return nil, ErrProformaNoEncontrada
	}

	actualizada.ID = id

	for i := range actualizada.Items {
		if actualizada.Items[i].ID == 0 {
			actualizada.Items[i].ID = s.nextItem
			s.nextItem++
		}
	}

	calcularTotales(&actualizada)
	s.proformas[id] = &actualizada

	return &actualizada, nil
}

// Calcular recalcula subtotales, impuestos y total de una proforma.
func (s *ProformaStorage) Calcular(id int) (*models.Proforma, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.proformas[id]
	if !ok {
		return nil, ErrProformaNoEncontrada
	}

	calcularTotales(p)

	copia := *p
	return &copia, nil
}

func calcularTotales(p *models.Proforma) {
	var subtotal float64

	for i := range p.Items {
		item := &p.Items[i]
		item.Subtotal = item.Cantidad * item.PrecioUnit
		subtotal += item.Subtotal
	}

	p.Subtotal = subtotal
	p.Impuestos = subtotal * tasaImpuesto
	p.Total = p.Subtotal + p.Impuestos
}
