
package storage

import (
	"errors"
	"sort"
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
)
//esto es para manejar los errores

var (
	ErrNotFound   = errors.New("recurso no encontrado") //esto es para manejar los errores en caso de que el recurso no exista
	ErrDuplicated = errors.New("nombre ya existe para esa unidad") //esto es para manejar los errores en caso de que el recurso ya exista
)


type Storage struct {
	mu         sync.RWMutex
	seq        int
	materiales map[int]*models.Material
	manoObras  map[int]*models.ManoObra
	equipos    map[int]*models.Equipo
	precios    map[int]*models.PrecioRecurso
	usuarios   map[int]*models.Usuario
}

func New() *Storage {
	return &Storage{
		materiales: make(map[int]*models.Material),
		manoObras:  make(map[int]*models.ManoObra),
		equipos:    make(map[int]*models.Equipo),
		precios:    make(map[int]*models.PrecioRecurso),
		usuarios:   make(map[int]*models.Usuario),
	}
}
//Este metodo se encarga de darle un id unico
func (s *Storage) nextID() int {
	s.seq++
	return s.seq
}

// ─────────────────────────────────────────────
// MATERIALES
// ─────────────────────────────────────────────

func (s *Storage) ListarMateriales() []*models.Material {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Material, 0, len(s.materiales))
	for _, m := range s.materiales {
		list = append(list, m)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (s *Storage) ObtenerMateriales(id int) (*models.Material, bool) { //este metodo se encarga de obtener el recurso
	s.mu.RLock() 
	defer s.mu.RUnlock()
	m, ok := s.materiales[id]
	if !ok {
		return nil, false
	}
	return m, true
}


func (s *Storage) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {  //este metodo se encarga de crear el recurso
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range s.materiales {
		if m.Nombre == in.Nombre && m.Unidad == in.Unidad { //si el recurso ya existe
			return nil, ErrDuplicated // se devuelve un error
		}
	}
	// Aquí se crea el recurso
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

func (s *Storage) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	mat, ok := s.materiales[id]
	if !ok {
		return nil, false
	}
	mat.Nombre = in.Nombre
	mat.Descripcion = in.Descripcion
	mat.Unidad = in.Unidad
	mat.PrecioReferencia = in.PrecioReferencia
	return mat, true
}

func (s *Storage) EliminarMateriales(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.materiales[id]; !ok {
		return false
	}
	delete(s.materiales, id)
	return true
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

func (s *Storage) ListarManoObra() []*models.ManoObra {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.ManoObra, 0, len(s.manoObras))
	for _, m := range s.manoObras {
		list = append(list, m)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (s *Storage) ObtenerManoObra(id int) (*models.ManoObra, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.manoObras[id]
	if !ok {
		return nil, false
	}
	return m, true
}

func (s *Storage) CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error) {
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

func (s *Storage) ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	mo, ok := s.manoObras[id]
	if !ok {
		return nil, false
	}
	mo.Descripcion = in.Descripcion
	mo.Categoria = in.Categoria
	mo.Unidad = in.Unidad
	mo.CostoReferencia = in.CostoReferencia
	return mo, true
}

func (s *Storage) EliminarManoObra(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.manoObras[id]; !ok {
		return false
	}
	delete(s.manoObras, id)
	return true
}

// ─────────────────────────────────────────────
// EQUIPOS
// ─────────────────────────────────────────────

func (s *Storage) ListarEquipos() []*models.Equipo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Equipo, 0, len(s.equipos))
	for _, e := range s.equipos {
		list = append(list, e)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (s *Storage) ObtenerEquipo(id int) (*models.Equipo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.equipos[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *Storage) CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error) {
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

func (s *Storage) ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error) {
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

func (s *Storage) EliminarEquipo(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.equipos[id]; !ok {
		return ErrNotFound
	}
	delete(s.equipos, id)
	return nil
}

// ─────────────────────────────────────────────
// PRECIOS
// ─────────────────────────────────────────────

func (s *Storage) ListarPrecios() []*models.PrecioRecurso {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.PrecioRecurso, 0, len(s.precios))
	for _, p := range s.precios {
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].FechaVigencia.Before(list[j].FechaVigencia)
	})
	return list
}

func (s *Storage) ObtenerPrecio(id int) (*models.PrecioRecurso, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.precios[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}
//existeRecurso es una función que se utiliza para verificar si un recurso existe en la base de datos.
func (s *Storage) ExisteRecurso(tipo string, id int) error {
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
	default:
		return ErrNotFound
	}
	return nil
}
//HistorialPrecios es una función que se utiliza para obtener el historial de precios de un recurso.
func (s *Storage) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.PrecioRecurso, 0)
	for _, p := range s.precios {
		if p.RecursoTipo == tipo && p.RecursoID == recursoID {
			list = append(list, p)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].FechaVigencia.Before(list[j].FechaVigencia)
	})
	return list
}
//PrecioVigente es una función que se utiliza para obtener el precio vigente de un recurso.
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

func (s *Storage) CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ExisteRecurso(in.RecursoTipo, in.RecursoID); err != nil {
		return nil, err
	}
	pr := &models.PrecioRecurso{
		ID:            s.nextID(),
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     in.RecursoID,
		Precio:        in.Precio,
		FechaVigencia: in.FechaVigencia,
		CreatedAt:     time.Now().UTC(),
	}
	s.precios[pr.ID] = pr
	return pr, nil
}

func (s *Storage) ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.precios[id]
	if !ok {
		return nil, ErrNotFound
	}
	p.Precio = in.Precio
	return p, nil
}

func (s *Storage) EliminarPrecio(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.precios[id]; !ok {
		return ErrNotFound
	}
	delete(s.precios, id)
	return nil
}

// ─────────────────────────────────────────────
// USUARIOS
// ─────────────────────────────────────────────

func (s *Storage) CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := &models.Usuario{
		ID:           s.nextID(),
		Email:        in.Email,
		PasswordHash: in.Password,
		CreatedAt:    time.Now().UTC(),
	}
	s.usuarios[u.ID] = u
	return u, nil
}

func (s *Storage) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.usuarios {
		if u.Email == email {
			return *u, true
		}
	}
	return models.Usuario{}, false
}

// ─────────────────────────────────────────────
// Memoria
// ─────────────────────────────────────────────

func (s *Storage) Seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()

	materiales := []struct{ nombre, descripcion, unidad, precio string }{
		{"Cemento Portland tipo I", "Saco de 50 kg para hormigón estructural", "unidad", "9.50"},
		{"Arena fina lavada", "Arena de río tamizada, libre de arcilla", "m³", "22.00"},
		{"Grava triturada 3/4\"", "Agregado grueso para hormigón", "m³", "28.50"},
		{"Varilla corrugada 12mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.15"},
		{"Varilla corrugada 16mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.12"},
		{"Bloque de hormigón 15x20x40", "Bloque vibrado para mampostería", "unidad", "0.65"},
		{"Ladrillo mambrón", "Ladrillo artesanal para fachada", "unidad", "0.28"},
		{"Tubo PVC presión 110mm", "Tubería PVC para agua potable, 6m", "unidad", "18.90"},
		{"Alambre de amarre #18", "Rollo de 25 kg", "kg", "1.80"},
		{"Pintura látex interior", "Pintura lavable blanca, 4 litros", "unidad", "12.80"},
		{"Cerámica piso 40x40", "Cerámica esmaltada para interior", "m²", "8.50"},
		{"Porcelanato 60x60", "Porcelanato rectificado mate", "m²", "18.00"},
		{"Impermeabilizante líquido", "Cementoso bicomponente, 20 kg", "unidad", "35.00"},
		{"Malla electrosoldada 15x15", "Acero para losas, rollo 6x2.4m", "unidad", "38.00"},
		{"Estuco listo para uso", "Pasta para empaste interior, 25 kg", "unidad", "7.20"},
	}
	for _, m := range materiales {
		id := s.nextID()
		s.materiales[id] = &models.Material{
			ID: id, Nombre: m.nombre, Descripcion: m.descripcion,
			Unidad: m.unidad, PrecioReferencia: decimal.RequireFromString(m.precio),
			CreatedAt: now,
		}
	}

	manoObras := []struct{ descripcion, categoria, unidad, costo string }{
		{"Maestro de obra general", "oficial", "día", "35.00"},
		{"Albañil - mampostería y enlucido", "oficial", "día", "28.00"},
		{"Fierrero - armado de acero", "oficial", "día", "30.00"},
		{"Carpintero de encofrado", "oficial", "día", "28.00"},
		{"Plomero instalaciones sanitarias", "oficial", "día", "32.00"},
		{"Electricista instalaciones", "oficial", "día", "32.00"},
		{"Pintor de interiores y exteriores", "oficial", "día", "26.00"},
		{"Ayudante de albañilería", "ayudante", "día", "18.00"},
		{"Ayudante de fierrero", "ayudante", "día", "18.00"},
		{"Peón de obra general", "ayudante", "día", "15.00"},
		{"Soldador estructura metálica", "especialista", "hora", "14.00"},
		{"Topógrafo de replanteo", "especialista", "día", "45.00"},
		{"Inspector de calidad hormigón", "especialista", "día", "50.00"},
	}
	for _, m := range manoObras {
		id := s.nextID()
		s.manoObras[id] = &models.ManoObra{
			ID: id, Descripcion: m.descripcion, Categoria: m.categoria,
			Unidad: m.unidad, CostoReferencia: decimal.RequireFromString(m.costo),
			CreatedAt: now,
		}
	}

	equipos := []struct {
		nombre, tipo, unidad, costoHora string
		disponible                      bool
	}{
		{"Concretera 1 saco eléctrica", "liviano", "hora", "8.50", true},
		{"Vibrador de hormigón 2\"", "liviano", "hora", "3.50", true},
		{"Amoladora angular 4.5\"", "liviano", "hora", "2.00", true},
		{"Compresor de aire 100 lt", "liviano", "hora", "4.00", true},
		{"Nivel láser rotativo", "liviano", "hora", "3.00", true},
		{"Excavadora sobre orugas CAT 320", "pesado", "hora", "85.00", true},
		{"Retroexcavadora JCB 3CX", "pesado", "hora", "65.00", true},
		{"Motoniveladora 140K", "pesado", "hora", "90.00", false},
		{"Compactadora tipo sapo", "pesado", "hora", "15.00", true},
		{"Volqueta 8 m³", "pesado", "hora", "55.00", true},
		{"Grúa torre 30m", "pesado", "hora", "120.00", false},
		{"Bomba de hormigón estacionaria", "pesado", "hora", "95.00", true},
	}
	for _, e := range equipos {
		id := s.nextID()
		s.equipos[id] = &models.Equipo{
			ID: id, Nombre: e.nombre, Tipo: e.tipo, Unidad: e.unidad,
			CostoHora:  decimal.RequireFromString(e.costoHora),
			Disponible: e.disponible, CreatedAt: now,
		}
	}

	hace30 := now.AddDate(0, -1, 0)
	hace60 := now.AddDate(0, -2, 0)

	type precioSeed struct {
		tipo   string
		id     int
		precio string
		fecha  time.Time
	}

	precios := []precioSeed{
		{"material", 1, "8.80", hace60},
		{"material", 1, "9.20", hace30},
		{"material", 1, "9.50", now},
		{"material", 2, "20.00", hace60},
		{"material", 2, "21.00", hace30},
		{"material", 2, "22.00", now},
		{"material", 4, "1.05", hace60},
		{"material", 4, "1.10", hace30},
		{"material", 4, "1.15", now},
		{"mano_obra", 16, "32.00", hace60},
		{"mano_obra", 16, "33.50", hace30},
		{"mano_obra", 16, "35.00", now},
		{"mano_obra", 17, "25.00", hace60},
		{"mano_obra", 17, "26.50", hace30},
		{"mano_obra", 17, "28.00", now},
		{"equipo", 34, "78.00", hace60},
		{"equipo", 34, "82.00", hace30},
		{"equipo", 34, "85.00", now},
		{"equipo", 35, "60.00", hace60},
		{"equipo", 35, "62.50", hace30},
		{"equipo", 35, "65.00", now},
	}

	for _, p := range precios {
		if err := s.ExisteRecurso(p.tipo, p.id); err != nil {
			continue
		}
		id := s.nextID()
		s.precios[id] = &models.PrecioRecurso{
			ID:            id,
			RecursoTipo:   p.tipo,
			RecursoID:     p.id,
			Precio:        decimal.RequireFromString(p.precio),
			FechaVigencia: p.fecha,
			CreatedAt:     p.fecha,
		}
	}
}
