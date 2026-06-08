package models
import(
	"errors"
	"time"
	"github.com/shopspring/decimal"
)

// ─────────────────────────────────────────────
// MATERIAL
// ─────────────────────────────────────────────

// UnidadesPermitidas es la lista de unidades válidas para materiales.
var UnidadesPermitidas = map[string]bool{
	"m²": true, "m³": true, "kg": true, "ton": true,
	"unidad": true, "jornal": true, "hora": true, "día": true,
	"ml": true, "lt": true, "gl": true,
}

type Material struct {
	ID               int             `json:"id"`
	Nombre           string          `json:"nombre"`
	Descripcion      string          `json:"descripcion"`
	Unidad           string          `json:"unidad"`
	PrecioReferencia decimal.Decimal `json:"precio_referencia"`
	CreatedAt        time.Time       `json:"created_at"`
}

// MaterialInput es el body que llega en POST / PUT.
type MaterialInput struct {
	Nombre           string          `json:"nombre"`
	Descripcion      string          `json:"descripcion"`
	Unidad           string          `json:"unidad"`
	PrecioReferencia decimal.Decimal `json:"precio_referencia"`
}

func (m MaterialInput) Validate() error {
	if m.Nombre == "" {
		return errors.New("campo 'nombre' requerido")
	}
	if m.Unidad == "" {
		return errors.New("campo 'unidad' requerido")
	}
	if !UnidadesPermitidas[m.Unidad] {
		return errors.New("unidad no permitida: use m², m³, kg, ton, unidad, jornal, hora, día, ml, lt, gl")
	}
	if m.PrecioReferencia.LessThanOrEqual(decimal.Zero) {
		return errors.New("precio_referencia debe ser mayor a 0")
	}
	return nil
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

var CategoriasPermitidas = map[string]bool{
	"oficial": true, "ayudante": true, "especialista": true,
}

var UnidadesManoObra = map[string]bool{
	"hora": true, "día": true, "jornal": true,
}

type ManoObra struct {
	ID              int             `json:"id"`
	Descripcion     string          `json:"descripcion"`
	Categoria       string          `json:"categoria"`
	Unidad          string          `json:"unidad"`
	CostoReferencia decimal.Decimal `json:"costo_referencia"`
	CreatedAt       time.Time       `json:"created_at"`
}

type ManoObraInput struct {
	Descripcion     string          `json:"descripcion"`
	Categoria       string          `json:"categoria"`
	Unidad          string          `json:"unidad"`
	CostoReferencia decimal.Decimal `json:"costo_referencia"`
}

func (m ManoObraInput) Validate() error {
	if m.Descripcion == "" {
		return errors.New("campo 'descripcion' requerido")
	}
	if m.Categoria == "" {
		return errors.New("campo 'categoria' requerido")
	}
	if !CategoriasPermitidas[m.Categoria] {
		return errors.New("categoria no permitida: use oficial, ayudante, especialista")
	}
	if m.Unidad == "" {
		return errors.New("campo 'unidad' requerido")
	}
	if !UnidadesManoObra[m.Unidad] {
		return errors.New("unidad no permitida para mano de obra: use hora, día, jornal")
	}
	if m.CostoReferencia.LessThanOrEqual(decimal.Zero) {
		return errors.New("costo_referencia debe ser mayor a 0")
	}
	return nil
}

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

var TiposEquipo = map[string]bool{
	"pesado": true, "liviano": true,
}

type Equipo struct {
	ID        int             `json:"id"`
	Nombre    string          `json:"nombre"`
	Tipo      string          `json:"tipo"`
	Unidad    string          `json:"unidad"`
	CostoHora decimal.Decimal `json:"costo_hora"`
	Disponible bool           `json:"disponible"`
	CreatedAt time.Time       `json:"created_at"`
}

type EquipoInput struct {
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Unidad     string          `json:"unidad"`
	CostoHora  decimal.Decimal `json:"costo_hora"`
	Disponible bool            `json:"disponible"`
}

func (e EquipoInput) Validate() error {
	if e.Nombre == "" {
		return errors.New("campo 'nombre' requerido")
	}
	if e.Tipo == "" {
		return errors.New("campo 'tipo' requerido")
	}
	if !TiposEquipo[e.Tipo] {
		return errors.New("tipo no permitido: use pesado, liviano")
	}
	if e.Unidad == "" {
		return errors.New("campo 'unidad' requerido")
	}
	if e.CostoHora.LessThanOrEqual(decimal.Zero) {
		return errors.New("costo_hora debe ser mayor a 0")
	}
	return nil
}

// DisponibilidadInput es el body del PATCH /equipos/:id/disponibilidad.
type DisponibilidadInput struct {
	Disponible bool `json:"disponible"`
}

// ─────────────────────────────────────────────
// PRECIO RECURSO 
// ─────────────────────────────────────────────

var RecursosTipos = map[string]bool{
	"material": true, "mano_obra": true, "equipo": true,
}

type PrecioRecurso struct {
	ID            int             `json:"id"`
	RecursoTipo   string          `json:"recurso_tipo"`
	RecursoID     int             `json:"recurso_id"`
	Precio        decimal.Decimal `json:"precio"`
	FechaVigencia time.Time       `json:"fecha_vigencia"`
	CreatedAt     time.Time       `json:"created_at"`
}

type PrecioRecursoInput struct {
	RecursoTipo   string          `json:"recurso_tipo"`
	RecursoID     int             `json:"recurso_id"`
	Precio        decimal.Decimal `json:"precio"`
	FechaVigencia time.Time       `json:"fecha_vigencia"`
}

func (p PrecioRecursoInput) Validate() error {
	if p.RecursoTipo == "" {
		return errors.New("campo 'recurso_tipo' requerido")
	}
	if !RecursosTipos[p.RecursoTipo] {
		return errors.New("recurso_tipo no válido: use material, mano_obra, equipo")
	}
	if p.RecursoID <= 0 {
		return errors.New("campo 'recurso_id' debe ser un entero positivo")
	}
	if p.Precio.LessThanOrEqual(decimal.Zero) {
		return errors.New("precio debe ser mayor a 0")
	}
	if p.FechaVigencia.IsZero() {
		return errors.New("campo 'fecha_vigencia' requerido")
	}
	return nil
}
