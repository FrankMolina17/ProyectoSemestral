package storage

import (
	"errors"
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

type ProformaStorage struct {
	mu            sync.Mutex
	proformas     map[int]models.Proforma
	items         map[int]models.ProformaItem
	notas         map[int]models.NotaProforma
	clientes      map[int]models.Cliente
	nextIDProf    int
	nextIDItem    int
	nextIDNota    int
	nextIDCliente int
}

func NuevoStorage() *ProformaStorage {
	return &ProformaStorage{
		proformas:     make(map[int]models.Proforma),
		items:         make(map[int]models.ProformaItem),
		notas:         make(map[int]models.NotaProforma),
		clientes:      make(map[int]models.Cliente),
		nextIDProf:    1,
		nextIDItem:    1,
		nextIDNota:    1,
		nextIDCliente: 1,
	}
}

// ── PROFORMAS ──

func (s *ProformaStorage) CrearProforma(p models.Proforma) models.Proforma {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextIDProf
	p.Estado = "borrador"
	p.CreadoEn = time.Now()
	p.Subtotal = 0
	p.Total = 0

	s.proformas[p.ID] = p
	s.nextIDProf++
	return p
}

func (s *ProformaStorage) ObtenerPorID(id int) (models.Proforma, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.proformas[id]
	if !ok {
		return models.Proforma{}, errors.New("proforma no encontrada")
	}
	return p, nil
}

func (s *ProformaStorage) ObtenerTodos() []models.Proforma {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]models.Proforma, 0, len(s.proformas))
	for _, p := range s.proformas {
		lista = append(lista, p)
	}
	return lista
}

func (s *ProformaStorage) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.proformas[id]
	if !ok {
		return models.Proforma{}, errors.New("proforma no encontrada")
	}

	p.Nombre = datos.Nombre
	p.PctGanancia = datos.PctGanancia
	p.PctImprevisto = datos.PctImprevisto
	p.ClienteID = datos.ClienteID

	s.proformas[id] = p
	return p, nil
}

func (s *ProformaStorage) EliminarProforma(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.proformas[id]
	if !ok {
		return errors.New("proforma no encontrada")
	}

	delete(s.proformas, id)
	return nil
}

func (s *ProformaStorage) AprobarProforma(id int) (models.Proforma, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.proformas[id]
	if !ok {
		return models.Proforma{}, errors.New("proforma no encontrada")
	}

	p.Estado = "aprobada"
	s.proformas[id] = p
	return p, nil
}

// ── ÍTEMS ──

func (s *ProformaStorage) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.proformas[item.ProformaID]
	if !ok {
		return models.ProformaItem{}, errors.New("proforma no encontrada")
	}

	item.ID = s.nextIDItem
	item.Subtotal = item.Cantidad * item.PrecioPromedio

	s.items[item.ID] = item
	s.nextIDItem++

	s.recalcularTotales(item.ProformaID)

	return item, nil
}

func (s *ProformaStorage) ObtenerItems(proformaID int) []models.ProformaItem {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]models.ProformaItem, 0)
	for _, item := range s.items {
		if item.ProformaID == proformaID {
			lista = append(lista, item)
		}
	}
	return lista
}

func (s *ProformaStorage) recalcularTotales(proformaID int) {
	p := s.proformas[proformaID]
	subtotal := 0.0

	for _, item := range s.items {
		if item.ProformaID == proformaID {
			subtotal += item.Subtotal
		}
	}

	p.Subtotal = subtotal
	p.Total = subtotal + (subtotal * p.PctGanancia) + (subtotal * p.PctImprevisto)
	s.proformas[proformaID] = p
}

// ── NOTAS ──

func (s *ProformaStorage) AgregarNota(nota models.NotaProforma) (models.NotaProforma, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.proformas[nota.ProformaID]
	if !ok {
		return models.NotaProforma{}, errors.New("proforma no encontrada")
	}

	nota.ID = s.nextIDNota
	nota.CreadoEn = time.Now()
	s.notas[nota.ID] = nota
	s.nextIDNota++
	return nota, nil
}

func (s *ProformaStorage) ObtenerNotas(proformaID int) []models.NotaProforma {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]models.NotaProforma, 0)
	for _, n := range s.notas {
		if n.ProformaID == proformaID {
			lista = append(lista, n)
		}
	}
	return lista
}

// ── CLIENTES ──

func (s *ProformaStorage) CrearCliente(c models.Cliente) models.Cliente {
	s.mu.Lock()
	defer s.mu.Unlock()

	c.ID = s.nextIDCliente
	s.clientes[c.ID] = c
	s.nextIDCliente++
	return c
}

func (s *ProformaStorage) ObtenerClientes() []models.Cliente {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]models.Cliente, 0, len(s.clientes))
	for _, c := range s.clientes {
		lista = append(lista, c)
	}
	return lista
}

func (s *ProformaStorage) ObtenerClientePorID(id int) (models.Cliente, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.clientes[id]
	if !ok {
		return models.Cliente{}, errors.New("cliente no encontrado")
	}
	return c, nil
}

func (s *ProformaStorage) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.clientes[id]
	if !ok {
		return models.Cliente{}, errors.New("cliente no encontrado")
	}

	c.Nombre = datos.Nombre
	c.Email = datos.Email
	c.Telefono = datos.Telefono
	c.Ruc = datos.Ruc

	s.clientes[id] = c
	return c, nil
}

func (s *ProformaStorage) EliminarCliente(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.clientes[id]
	if !ok {
		return errors.New("cliente no encontrado")
	}

	delete(s.clientes, id)
	return nil
}