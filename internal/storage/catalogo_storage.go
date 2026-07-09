package storage

import (
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

type CatalogoStorage struct {
	mu          sync.Mutex
	materiales  map[int]models.Material
	manoObras   map[int]models.ManoObra
	equipos     map[int]models.Equipo
	precios     map[int]models.Precio
	usuarios    map[int]models.UsuarioCatalogo
	nextMatID   int
	nextMobID   int
	nextEqID    int
	nextPreID   int
	nextUserID  int
}

func NuevoCatalogoStorage() *CatalogoStorage {
	return &CatalogoStorage{
		materiales: make(map[int]models.Material),
		manoObras:  make(map[int]models.ManoObra),
		equipos:    make(map[int]models.Equipo),
		precios:    make(map[int]models.Precio),
		usuarios:   make(map[int]models.UsuarioCatalogo),
		nextMatID:  1,
		nextMobID:  1,
		nextEqID:   1,
		nextPreID:  1,
		nextUserID: 1,
	}
}

func (s *CatalogoStorage) Seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.materiales[1] = models.Material{ID: 1, Nombre: "Cemento Portland", Unidad: "bolsa", Cantidad: 500, CostoUnit: 25.50, Proveedor: "Constructora S.A."}
	s.materiales[2] = models.Material{ID: 2, Nombre: "Arena Fina", Unidad: "m³", Cantidad: 100, CostoUnit: 120.00, Proveedor: "Arenera del Sur"}
	s.materiales[3] = models.Material{ID: 3, Nombre: "Ladrillo Hueco", Unidad: "unidad", Cantidad: 10000, CostoUnit: 3.75, Proveedor: "Ladrillera Norte"}
	s.nextMatID = 4

	s.manoObras[1] = models.ManoObra{ID: 1, Nombre: "Albañil", Tipo: "general", CostoPorHora: 45.00}
	s.manoObras[2] = models.ManoObra{ID: 2, Nombre: "Electricista", Tipo: "especializado", CostoPorHora: 60.00}
	s.manoObras[3] = models.ManoObra{ID: 3, Nombre: "Gasfitero", Tipo: "especializado", CostoPorHora: 55.00}
	s.nextMobID = 4

	s.equipos[1] = models.Equipo{ID: 1, Nombre: "Mezcladora de Concreto", Modelo: "MZ-300", CostoPorHora: 120.00}
	s.equipos[2] = models.Equipo{ID: 2, Nombre: "Compresor de Aire", Modelo: "CA-150", CostoPorHora: 85.00}
	s.equipos[3] = models.Equipo{ID: 3, Nombre: "Vibrador de Concreto", Modelo: "VC-50", CostoPorHora: 40.00}
	s.nextEqID = 4

	ahora := time.Now()
	s.precios[1] = models.Precio{ID: 1, Tipo: "material", RecursoID: 1, Monto: 25.50, Vigente: true, CreadoEn: ahora}
	s.precios[2] = models.Precio{ID: 2, Tipo: "mano_obra", RecursoID: 1, Monto: 45.00, Vigente: true, CreadoEn: ahora}
	s.precios[3] = models.Precio{ID: 3, Tipo: "equipo", RecursoID: 1, Monto: 120.00, Vigente: true, CreadoEn: ahora}
	s.nextPreID = 4

	s.usuarios[1] = models.UsuarioCatalogo{ID: 1, Email: "admin@catalogo.com", Password: "admin123", Nombre: "Admin Catálogo"}
	s.nextUserID = 2
}

func (s *CatalogoStorage) ListarMateriales() []models.Material {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]models.Material, 0, len(s.materiales))
	for _, m := range s.materiales {
		res = append(res, m)
	}
	return res
}

func (s *CatalogoStorage) ObtenerMaterial(id int) (models.Material, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.materiales[id]
	return m, ok
}

func (s *CatalogoStorage) CrearMaterial(m models.Material) models.Material {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextMatID
	s.nextMatID++
	s.materiales[m.ID] = m
	return m
}

func (s *CatalogoStorage) ActualizarMaterial(m models.Material) (models.Material, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.materiales[m.ID]
	if !ok {
		return models.Material{}, false
	}
	s.materiales[m.ID] = m
	return m, true
}

func (s *CatalogoStorage) BorrarMaterial(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.materiales[id]
	if !ok {
		return false
	}
	delete(s.materiales, id)
	return true
}

func (s *CatalogoStorage) ListarManoObras() []models.ManoObra {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]models.ManoObra, 0, len(s.manoObras))
	for _, m := range s.manoObras {
		res = append(res, m)
	}
	return res
}

func (s *CatalogoStorage) ObtenerManoObra(id int) (models.ManoObra, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.manoObras[id]
	return m, ok
}

func (s *CatalogoStorage) CrearManoObra(m models.ManoObra) models.ManoObra {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextMobID
	s.nextMobID++
	s.manoObras[m.ID] = m
	return m
}

func (s *CatalogoStorage) ActualizarManoObra(m models.ManoObra) (models.ManoObra, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.manoObras[m.ID]
	if !ok {
		return models.ManoObra{}, false
	}
	s.manoObras[m.ID] = m
	return m, true
}

func (s *CatalogoStorage) BorrarManoObra(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.manoObras[id]
	if !ok {
		return false
	}
	delete(s.manoObras, id)
	return true
}

func (s *CatalogoStorage) ListarEquipos() []models.Equipo {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]models.Equipo, 0, len(s.equipos))
	for _, e := range s.equipos {
		res = append(res, e)
	}
	return res
}

func (s *CatalogoStorage) ObtenerEquipo(id int) (models.Equipo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.equipos[id]
	return e, ok
}

func (s *CatalogoStorage) CrearEquipo(e models.Equipo) models.Equipo {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = s.nextEqID
	s.nextEqID++
	s.equipos[e.ID] = e
	return e
}

func (s *CatalogoStorage) ActualizarEquipo(e models.Equipo) (models.Equipo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.equipos[e.ID]
	if !ok {
		return models.Equipo{}, false
	}
	s.equipos[e.ID] = e
	return e, true
}

func (s *CatalogoStorage) BorrarEquipo(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.equipos[id]
	if !ok {
		return false
	}
	delete(s.equipos, id)
	return true
}

func (s *CatalogoStorage) ListarPrecios() []models.Precio {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]models.Precio, 0, len(s.precios))
	for _, p := range s.precios {
		res = append(res, p)
	}
	return res
}

func (s *CatalogoStorage) ObtenerPrecio(id int) (models.Precio, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.precios[id]
	return p, ok
}

func (s *CatalogoStorage) CrearPrecio(p models.Precio) models.Precio {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = s.nextPreID
	s.nextPreID++
	p.CreadoEn = time.Now()
	s.precios[p.ID] = p
	return p
}

func (s *CatalogoStorage) ActualizarPrecio(p models.Precio) (models.Precio, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.precios[p.ID]
	if !ok {
		return models.Precio{}, false
	}
	s.precios[p.ID] = p
	return p, true
}

func (s *CatalogoStorage) BorrarPrecio(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.precios[id]
	if !ok {
		return false
	}
	delete(s.precios, id)
	return true
}

func (s *CatalogoStorage) PreciosPorRecurso(tipo string, recursoID int) []models.Precio {
	s.mu.Lock()
	defer s.mu.Unlock()
	var res []models.Precio
	for _, p := range s.precios {
		if p.Tipo == tipo && p.RecursoID == recursoID {
			res = append(res, p)
		}
	}
	return res
}

func (s *CatalogoStorage) PrecioVigente(tipo string, recursoID int) (models.Precio, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.precios {
		if p.Tipo == tipo && p.RecursoID == recursoID && p.Vigente {
			return p, true
		}
	}
	return models.Precio{}, false
}

func (s *CatalogoStorage) BuscarUsuarioCatalogo(email, password string) (models.UsuarioCatalogo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.usuarios {
		if u.Email == email && u.Password == password {
			return u, true
		}
	}
	return models.UsuarioCatalogo{}, false
}
