package storage
import (
	"errors"
	"sort"
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

// ─────────────────────────────────────────────
// Errores de dominio
// ─────────────────────────────────────────────

var (
	ErrNotFound   = errors.New("recurso no encontrado")
	ErrDuplicated = errors.New("nombre ya existe para esa unidad")
)


//Storage es un repositorio en memoria

type Storage struct {
	mu  sync.RWMutex
	seq int

	materiales map[int]*models.Material
	manoObras  map[int]*models.ManoObra
	equipos    map[int]*models.Equipo
	precios    map[int]*models.PrecioRecurso
}
func New() *Storage {
	return &Storage{
		materiales: make(map[int]*models.Material),
		manoObras:  make(map[int]*models.ManoObra),
		equipos:    make(map[int]*models.Equipo),
		precios:    make(map[int]*models.PrecioRecurso),
	}
}

func (s *Storage) nextID() int {
	s.seq++
	return s.seq
}

// ─────────────────────────────────────────────
// MATERIALES
// ─────────────────────────────────────────────

func (s *Storage) ListMateriales(nombre, unidad string) []*models.Material {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*models.Material
	for _, m := range s.materiales {
		if nombre != "" && m.Nombre != nombre {
			continue
		}
		if unidad != "" && m.Unidad != unidad {
			continue
		}
		list = append(list, m)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

//Obtener Materiales por id
func (s *Storage) ObtenerMaterialporID(id int) (*models.Material, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.materiales[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}
func (s *Storage) CrearMaterial(in models.MaterialInput) (*models.Material, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// unicidad por nombre+unidad
	for _, m := range s.materiales {
		if m.Nombre == in.Nombre && m.Unidad == in.Unidad {
			return nil, ErrDuplicated
		}
	}

	mat := &models.Material{
		ID:               s.nextID(),
		Nombre:           in.Nombre,
		Descripcion:      in.Descripcion,
		Unidad:           in.Unidad,
		PrecioReferencia: in.PrecioReferencia,
		CreatedAt:        time.Now().UTC(),
	}
	s.materiales[mat.ID] = mat
	return mat, nil
}


func (s *Storage) ActualizarMaterial(id int, in models.MaterialInput) (*models.Material, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mat, ok := s.materiales[id]
	if !ok {
		return nil, ErrNotFound
	}
	mat.Nombre = in.Nombre
	mat.Descripcion = in.Descripcion
	mat.Unidad = in.Unidad
	mat.PrecioReferencia = in.PrecioReferencia
	return mat, nil
}

//borrar
func (s *Storage) BorrarMaterial(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.materiales[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.materiales, id)
	return nil
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

//listar mano de obra

func (s *Storage) ListarManoObra(categoria string) []*models.ManoObra {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*models.ManoObra
	for _, m := range s.manoObras {
		if categoria != "" && m.Categoria != categoria {
			continue
		}
		list = append(list, m)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

//obtener mano de obra por id
func (s *Storage) ObtenerManoObraPorID(id int) (*models.ManoObra, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.manoObras[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

//crear mano de obra
func (s *Storage) CrearManoObra(in models.ManoObraInput) (*models.ManoObra, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mo := &models.ManoObra{
		ID:              s.nextID(),
		Descripcion:     in.Descripcion,
		Categoria:       in.Categoria,
		Unidad:          in.Unidad,
		CostoReferencia: in.CostoReferencia,
		CreatedAt:       time.Now().UTC(),
	}
	s.manoObras[mo.ID] = mo
	return mo, nil


}

//acualizar
func (s *Storage) ActualizarManoObra(id int, in models.ManoObraInput) (*models.ManoObra, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	mo, ok := s.manoObras[id]
	if !ok {
		return nil, ErrNotFound
	}
	mo.Descripcion = in.Descripcion
	mo.Categoria = in.Categoria
	mo.Unidad = in.Unidad
	mo.CostoReferencia = in.CostoReferencia
	return mo, nil
}

//borrar
func (s *Storage) BorrarManoObra(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.manoObras[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.manoObras, id)
	return nil
}

// ─────────────────────────────────────────────
// EQUIPOS
// ─────────────────────────────────────────────
func (s *Storage) ListarEquipos(disponible *bool, tipo string) []*models.Equipo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*models.Equipo
	for _, e := range s.equipos {
		if disponible != nil && e.Disponible != *disponible {
			continue
		}
		if tipo != "" && e.Tipo != tipo {
			continue
		}
		list = append(list, e)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

//obtener equipo por id
func (s *Storage) ObtenerEquipoPorID(id int) (*models.Equipo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.equipos[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *Storage) CrearEquipo(in models.EquipoInput) (*models.Equipo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	eq := &models.Equipo{
		ID:         s.nextID(),
		Nombre:     in.Nombre,
		Tipo:       in.Tipo,
		Unidad:     in.Unidad,
		CostoHora:  in.CostoHora,
		Disponible: in.Disponible,
		CreatedAt:  time.Now().UTC(),
	}
	s.equipos[eq.ID] = eq
	return eq, nil
}

func (s *Storage) ActualizarEquipo(id int, in models.EquipoInput) (*models.Equipo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	eq, ok := s.equipos[id]
	if !ok {
		return nil, ErrNotFound
	}
	eq.Nombre = in.Nombre
	eq.Tipo = in.Tipo
	eq.Unidad = in.Unidad
	eq.CostoHora = in.CostoHora
	eq.Disponible = in.Disponible
	return eq, nil
}

func (s *Storage) BorrarEquipo(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.equipos[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.equipos, id)
	return nil
}

func (s *Storage) Disponibilidad(id int, disponible bool) (*models.Equipo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	eq, ok := s.equipos[id]
	if !ok {
		return nil, ErrNotFound
	}
	eq.Disponible = disponible
	return eq, nil
}

// ─────────────────────────────────────────────
// PRECIOS
// ─────────────────────────────────────────────
// existeRecurso verifica que el recurso con el id dado exista.

func (s *Storage) existeRecurso(tipo string, id int) error {
	switch tipo {
	case "material":
		if _, ok := s.materiales[id]; !ok {
			return ErrNotFound
		}
	case "mano_obra":
		if _, ok := s.manoObras[id]; !ok {
			return ErrNotFound
		}
	case "equipo":
		if _, ok := s.equipos[id]; !ok {
			return ErrNotFound
		}
	}
	return nil
}

func (s *Storage) CreatePrecio(in models.PrecioRecursoInput) (*models.PrecioRecurso, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// verificar que el recurso referenciado exista
	if err := s.existeRecurso(in.RecursoTipo, in.RecursoID); err != nil {
		return nil, err
	}

	pr := &models.PrecioRecurso{
		ID:            s.nextID(),
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     in.RecursoID,
		Precio:        in.Precio,
		FechaVigencia: in.FechaVigencia.UTC(),
		CreatedAt:     time.Now().UTC(),
	}
	s.precios[pr.ID] = pr
	return pr, nil
}

func (s *Storage) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*models.PrecioRecurso
	for _, p := range s.precios {
		if p.RecursoTipo == tipo && p.RecursoID == recursoID {
			list = append(list, p)
		}
	}
	// ordenar de más antiguo a más nuevo
	sort.Slice(list, func(i, j int) bool {
		return list[i].FechaVigencia.Before(list[j].FechaVigencia)
	})
	return list
}


// PrecioVigente devuelve el precio con la FechaVigencia más reciente
// que sea ≤ now (el precio "en efecto" hoy).
func (s *Storage) PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now().UTC()
	var vigente *models.PrecioRecurso
	for _, p := range s.precios {
		if p.RecursoTipo != tipo || p.RecursoID != recursoID {
			continue
		}
		if p.FechaVigencia.After(now) {
			continue
		}
		if vigente == nil || p.FechaVigencia.After(vigente.FechaVigencia) {
			vigente = p
		}
	}
	if vigente == nil {
		return nil, ErrNotFound
	}
	return vigente, nil
}

