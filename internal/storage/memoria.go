package storage

import (
	"errors"
	"sync"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

var (
	ErrNotFound   = errors.New("recurso no encontrado")
	ErrDuplicated = errors.New("nombre ya existe para esa unidad")
)

type Memoria struct {
	incidencias      []models.Incidencia
	nextIncidenciaID int

	mu sync.Mutex
}

// ListarProductos devuelve todos los productos en memoria.
func (m *Memoria) ListarIncidencias() []models.Incidencia {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Incidencia, len(m.incidencias))
	copy(copia, m.incidencias)
	return copia
}

// BuscarProductoPorID devuelve el producto con el ID dado (patrón comma-ok).
func (m *Memoria) BuscarIncidenciaPorID(id int) (models.Incidencia, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.incidencias {
		if p.ID == id {
			return p, true
		}
	}
	return models.Incidencia{}, false
}

func (m *Memoria) BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.incidencias {
		if p.ID == id && p.EntidadTipo == tipo {
			return p, true
		}
	}
	return models.Incidencia{}, false
}

// CrearProducto agrega un producto nuevo y devuelve el producto con ID asignado.
func (m *Memoria) CrearIncidencia(p models.Incidencia) models.Incidencia {
	m.mu.Lock()
	defer m.mu.Unlock()

	p.ID = m.nextIncidenciaID
	m.nextIncidenciaID++
	m.incidencias = append(m.incidencias, p)
	return p
}

// ActualizarProducto reemplaza el producto con el ID dado.
func (m *Memoria) ActualizarIncidencia(id int, datos models.Incidencia) (models.Incidencia, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.incidencias {
		if p.ID == id {
			datos.ID = id
			m.incidencias[i] = datos
			return datos, true
		}
	}
	return models.Incidencia{}, false
}

// BorrarProducto elimina el producto con el ID dado.
func (m *Memoria) BorrarIncidencia(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.incidencias {
		if p.ID == id {
			m.incidencias = append(m.incidencias[:i], m.incidencias[i+1:]...)
			return true
		}
	}
	return false
}
