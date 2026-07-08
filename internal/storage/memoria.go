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
	nextUsuarioID    int
	obras            []models.Obra
	usuarios         []models.Usuario
	nextObraID       int
	mu               sync.Mutex
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

// =========================================================
// MÉTODOS PARA OBRAS (NUEVO)
// =========================================================

func (m *Memoria) ListarObras() []models.Obra {
	m.mu.Lock()
	defer m.mu.Unlock()

	copia := make([]models.Obra, len(m.obras))
	copy(copia, m.obras)
	return copia
}

func (m *Memoria) BuscarObraPorID(id int) (models.Obra, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, o := range m.obras {
		if o.ID == id {
			return o, true
		}
	}
	return models.Obra{}, false
}

func (m *Memoria) CrearObra(o models.Obra) models.Obra {
	m.mu.Lock()
	defer m.mu.Unlock()

	o.ID = m.nextObraID
	m.nextObraID++
	m.obras = append(m.obras, o)
	return o
}

func (m *Memoria) ActualizarObra(id int, datos models.Obra) (models.Obra, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, o := range m.obras {
		if o.ID == id {
			datos.ID = id
			m.obras[i] = datos
			return datos, true
		}
	}
	return models.Obra{}, false
}

func (m *Memoria) BorrarObra(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, o := range m.obras {
		if o.ID == id {
			m.obras = append(m.obras[:i], m.obras[i+1:]...)
			return true
		}
	}
	return false
}

func (m *Memoria) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

func (m *Memoria) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	u.ID = m.nextUsuarioID
	m.nextUsuarioID++
	m.usuarios = append(m.usuarios, u)
	return u, nil
}
