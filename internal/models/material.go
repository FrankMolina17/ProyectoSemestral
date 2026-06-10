package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// ─────────────────────────────────────────────
// MATERIAL
// ─────────────────────────────────────────────

type Material struct {
	ID               int             `json:"id"`
	Nombre           string          `json:"nombre"`
	Descripcion      string          `json:"descripcion"`
	Unidad           string          `json:"unidad"`
	PrecioReferencia decimal.Decimal `json:"precio_referencia"`
	CreatedAt        time.Time       `json:"created_at"`
}

// MaterialInput es una estructura que representa la información de entrada para un material. 
// Contiene los campos Nombre, Descripcion, Unidad y PrecioReferencia.
//si la unidad no es permitida retorna un error
type EntradaMaterial struct {
	Nombre           string          `json:"nombre"`
	Descripcion      string          `json:"descripcion"`
	Unidad           string          `json:"unidad"`
	PrecioReferencia decimal.Decimal `json:"precio_referencia"`
}

var unidadesPermitidas = map[string]bool{
	"m²": true, "m³": true, "kg": true, "ton": true,
	"unidad": true, "jornal": true, "hora": true, "día": true,
	"ml": true, "lt": true, "gl": true,
}

func (m EntradaMaterial) ValidarMaterial() error {
	if m.Nombre == "" {
		return errors.New("campo 'nombre' requerido")
	}
	if m.Unidad == "" {
		return errors.New("campo 'unidad' requerido")
	}
	if !unidadesPermitidas[m.Unidad] {
		return errors.New("unidad no permitida")
	}
	if m.PrecioReferencia.LessThanOrEqual(decimal.Zero) {
		return errors.New("precio_referencia debe ser mayor a 0")
	}
	return nil
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

type ManoObra struct {
	ID              int             `json:"id"`
	Descripcion     string          `json:"descripcion"`
	Categoria       string          `json:"categoria"`
	Unidad          string          `json:"unidad"`
	CostoReferencia decimal.Decimal `json:"costo_referencia"`
	CreatedAt       time.Time       `json:"created_at"`
}
type EntradaManoObra struct {
	Descripcion     string          `json:"descripcion"`
	Categoria       string          `json:"categoria"`
	Unidad          string          `json:"unidad"`
	CostoReferencia decimal.Decimal `json:"costo_referencia"`
}

var categoriasPermitidas = map[string]bool{
	"oficial": true, "ayudante": true, "especialista": true,
}

var unidadesManoObra = map[string]bool{
	"hora": true, "día": true, "jornal": true,
}

func (m EntradaManoObra) ValidarManoObra() error {
	if m.Descripcion == "" {
		return errors.New("campo 'descripcion' requerido")
	}
	if !categoriasPermitidas[m.Categoria] {
		return errors.New("categoria debe ser: oficial, ayudante, especialista")
	}
	if !unidadesManoObra[m.Unidad] {
		return errors.New("unidad debe ser: hora, día, jornal")
	}
	if m.CostoReferencia.LessThanOrEqual(decimal.Zero) {
		return errors.New("costo_referencia debe ser mayor a 0")
	}
	return nil
}

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

type Equipo struct {
	ID         int             `json:"id"`
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Unidad     string          `json:"unidad"`
	CostoHora  decimal.Decimal `json:"costo_hora"`
	Disponible bool            `json:"disponible"`
	CreatedAt  time.Time       `json:"created_at"`
}

type EntradaEquipo struct {
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Unidad     string          `json:"unidad"`
	CostoHora  decimal.Decimal `json:"costo_hora"`
	Disponible bool            `json:"disponible"`
}

var tiposEquipo = map[string]bool{
	"pesado": true, "liviano": true,
}

func (e EntradaEquipo) ValidarEquipo() error {
	if e.Nombre == "" {
		return errors.New("campo 'nombre' requerido")
	}
	if !tiposEquipo[e.Tipo] {
		return errors.New("tipo debe ser: pesado, liviano")
	}
	if e.Unidad == "" {
		return errors.New("campo 'unidad' requerido")
	}
	if e.CostoHora.LessThanOrEqual(decimal.Zero) {
		return errors.New("costo_hora debe ser mayor a 0")
	}
	return nil
}

// ─────────────────────────────────────────────
// PRECIO RECURSO
// ─────────────────────────────────────────────

type PrecioRecurso struct {
	ID            int             `json:"id"`
	RecursoTipo   string          `json:"recurso_tipo"`
	RecursoID     int             `json:"recurso_id"`
	Precio        decimal.Decimal `json:"precio"`
	FechaVigencia time.Time       `json:"fecha_vigencia"`
	CreatedAt     time.Time       `json:"created_at"`
}

type EntradaPrecioRecurso struct {
	RecursoTipo   string          `json:"recurso_tipo"`
	RecursoID     int             `json:"recurso_id"`
	Precio        decimal.Decimal `json:"precio"`
	FechaVigencia time.Time       `json:"fecha_vigencia"`
}

var RecursosTipos = map[string]bool{
	"material": true, "mano_obra": true, "equipo": true,
}

func (p EntradaPrecioRecurso) ValidarPrecio() error {
	if !RecursosTipos[p.RecursoTipo] {
		return errors.New("recurso_tipo debe ser: material, mano_obra, equipo")
	}
	if p.RecursoID <= 0 {
		return errors.New("recurso_id debe ser un entero positivo")
	}
	if p.Precio.LessThanOrEqual(decimal.Zero) {
		return errors.New("precio debe ser mayor a 0")
	}
	if p.FechaVigencia.IsZero() {
		return errors.New("campo 'fecha_vigencia' requerido")
	}
	return nil
}
