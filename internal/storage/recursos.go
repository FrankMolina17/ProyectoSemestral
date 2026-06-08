package storage
import(
	"time"

	"github.com/shopspring/decimal"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)



// Seed carga datos de prueba realistas en el store.
// Se llama una sola vez desde main.go al iniciar el servidor.
func (s *Storage) Seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	hace30 := now.AddDate(0, -1, 0)
	hace60 := now.AddDate(0, -2, 0)

	// ─────────────────────────────────────────────
	// MATERIALES
	// ─────────────────────────────────────────────

	materiales := []struct {
		nombre      string
		descripcion string
		unidad      string
		precio      string
	}{
		{"Cemento Portland tipo I", "Saco de 50 kg para hormigón estructural", "unidad", "9.50"},
		{"Arena fina lavada", "Arena de río tamizada, libre de arcilla", "m³", "22.00"},
		{"Grava triturada 3/4\"", "Agregado grueso para hormigón", "m³", "28.50"},
		{"Varilla corrugada 12mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.15"},
		{"Varilla corrugada 16mm", "Acero de refuerzo ASTM A615 Gr.60", "kg", "1.12"},
		{"Bloque de hormigón 15x20x40", "Bloque vibrado para mampostería", "unidad", "0.65"},
		{"Ladrillo mambrón", "Ladrillo artesanal para fachada", "unidad", "0.28"},
		{"Tubo PVC presión 110mm", "Tubería PVC para agua potable, 6m", "unidad", "18.90"},
		{"Tubo PVC sanitario 110mm", "Tubería sanitaria, 3m", "unidad", "8.40"},
		{"Alambre de amarre #18", "Rollo de 25 kg para amarre de acero", "kg", "1.80"},
		{"Clavos 2\"", "Clavos de acero galvanizado", "kg", "1.20"},
		{"Madera de encofrado 1x10", "Tabla de laurel para encofrado", "unidad", "4.50"},
		{"Puntales metálicos 3m", "Puntal telescópico para losas", "unidad", "45.00"},
		{"Pintura látex interior", "Pintura lavable blanca, 4 litros", "unidad", "12.80"},
		{"Cerámica piso 40x40", "Cerámica esmaltada para interior", "m²", "8.50"},
		{"Porcelanato 60x60", "Porcelanato rectificado mate", "m²", "18.00"},
		{"Vidrio claro 6mm", "Vidrio flotado para ventanas", "m²", "22.50"},
		{"Impermeabilizante líquido", "Impermeabilizante cementoso bicomponente, 20 kg", "unidad", "35.00"},
		{"Malla electrosoldada 15x15", "Malla de acero para losas, rollo 6x2.4m", "unidad", "38.00"},
		{"Estuco listo para uso", "Pasta para empaste interior, saco 25 kg", "unidad", "7.20"},
	}

	for _, m := range materiales {
		id := s.nextID()
		s.materiales[id] = &models.Material{
			ID:               id,
			Nombre:           m.nombre,
			Descripcion:      m.descripcion,
			Unidad:           m.unidad,
			PrecioReferencia: decimal.RequireFromString(m.precio),
			CreatedAt:        now,
		}
	}

	// ─────────────────────────────────────────────
	// MANO DE OBRA
	// ─────────────────────────────────────────────

	manoObras := []struct {
		descripcion string
		categoria   string
		unidad      string
		costo       string
	}{
		{"Maestro de obra general", "oficial", "día", "35.00"},
		{"Albañil - mampostería y enlucido", "oficial", "día", "28.00"},
		{"Fierrero - armado de acero", "oficial", "día", "30.00"},
		{"Carpintero de encofrado", "oficial", "día", "28.00"},
		{"Plomero instalaciones sanitarias", "oficial", "día", "32.00"},
		{"Electricista instalaciones", "oficial", "día", "32.00"},
		{"Operador de maquinaria pesada", "oficial", "hora", "12.00"},
		{"Pintor de interiores y exteriores", "oficial", "día", "26.00"},
		{"Ayudante de albañilería", "ayudante", "día", "18.00"},
		{"Ayudante de fierrero", "ayudante", "día", "18.00"},
		{"Ayudante de plomería", "ayudante", "día", "16.00"},
		{"Ayudante de electricista", "ayudante", "día", "16.00"},
		{"Peón de obra general", "ayudante", "día", "15.00"},
		{"Colocador de cerámica y porcelanato", "oficial", "día", "28.00"},
		{"Soldador estructura metálica", "especialista", "hora", "14.00"},
		{"Topógrafo de replanteo", "especialista", "día", "45.00"},
		{"Inspector de calidad hormigón", "especialista", "día", "50.00"},
	}

	for _, mo := range manoObras {
		id := s.nextID()
		s.manoObras[id] = &models.ManoObra{
			ID:              id,
			Descripcion:     mo.descripcion,
			Categoria:       mo.categoria,
			Unidad:          mo.unidad,
			CostoReferencia: decimal.RequireFromString(mo.costo),
			CreatedAt:       now,
		}
	}

	// ─────────────────────────────────────────────
	// EQUIPOS
	// ─────────────────────────────────────────────

	equipos := []struct {
		nombre     string
		tipo       string
		unidad     string
		costoHora  string
		disponible bool
	}{
		{"Concretera 1 saco eléctrica", "liviano", "hora", "8.50", true},
		{"Concretera 2 sacos a gasolina", "liviano", "hora", "12.00", true},
		{"Vibrador de hormigón 2\"", "liviano", "hora", "3.50", true},
		{"Andamio tubular (por módulo)", "liviano", "hora", "0.80", true},
		{"Amoladora angular 4.5\"", "liviano", "hora", "2.00", true},
		{"Taladro percutor 1/2\"", "liviano", "hora", "1.50", true},
		{"Compresor de aire 100 lt", "liviano", "hora", "4.00", true},
		{"Nivel láser rotativo", "liviano", "hora", "3.00", true},
		{"Excavadora sobre orugas CAT 320", "pesado", "hora", "85.00", true},
		{"Retroexcavadora JCB 3CX", "pesado", "hora", "65.00", true},
		{"Motoniveladora 140K", "pesado", "hora", "90.00", false},
		{"Compactadora tipo sapo", "pesado", "hora", "15.00", true},
		{"Rodillo vibratorio 1.5 ton", "pesado", "hora", "45.00", true},
		{"Volqueta 8 m³", "pesado", "hora", "55.00", true},
		{"Grúa torre 30m", "pesado", "hora", "120.00", false},
		{"Bomba de hormigón estacionaria", "pesado", "hora", "95.00", true},
		{"Generador eléctrico 50 KVA", "liviano", "hora", "18.00", true},
	}

	for _, e := range equipos {
		id := s.nextID()
		s.equipos[id] = &models.Equipo{
			ID:         id,
			Nombre:     e.nombre,
			Tipo:       e.tipo,
			Unidad:     e.unidad,
			CostoHora:  decimal.RequireFromString(e.costoHora),
			Disponible: e.disponible,
			CreatedAt:  now,
		}
	}

	// ─────────────────────────────────────────────
	// PRECIOS HISTÓRICOS
	// Solo para los primeros recursos de cada tipo,
	// para ilustrar el historial de variación de precios.
	// ─────────────────────────────────────────────

	type precioSeed struct {
		tipo   string
		id     int
		precio string
		fecha  time.Time
	}

	// IDs asignados: materiales 1-20, manoObras 21-37, equipos 38-54
	// Precio hace 60 días → hace 30 días → hoy (ya en PrecioReferencia)
	precios := []precioSeed{
		// Cemento Portland (id=1)
		{"material", 1, "8.80", hace60},
		{"material", 1, "9.20", hace30},
		{"material", 1, "9.50", now},
		// Arena fina (id=2)
		{"material", 2, "20.00", hace60},
		{"material", 2, "21.00", hace30},
		{"material", 2, "22.00", now},
		// Varilla 12mm (id=4)
		{"material", 4, "1.05", hace60},
		{"material", 4, "1.10", hace30},
		{"material", 4, "1.15", now},
		// Maestro de obra (id=21)
		{"mano_obra", 21, "32.00", hace60},
		{"mano_obra", 21, "33.50", hace30},
		{"mano_obra", 21, "35.00", now},
		// Albañil (id=22)
		{"mano_obra", 22, "25.00", hace60},
		{"mano_obra", 22, "26.50", hace30},
		{"mano_obra", 22, "28.00", now},
		// Excavadora CAT 320 (id=46)
		{"equipo", 46, "78.00", hace60},
		{"equipo", 46, "82.00", hace30},
		{"equipo", 46, "85.00", now},
		// Retroexcavadora JCB (id=47)
		{"equipo", 47, "60.00", hace60},
		{"equipo", 47, "62.50", hace30},
		{"equipo", 47, "65.00", now},
	}

	for _, p := range precios {
		// verificar que el recurso exista antes de insertar
		if err := s.existeRecurso(p.tipo, p.id); err != nil {
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
