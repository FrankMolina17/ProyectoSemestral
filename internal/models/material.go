package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Material struct {
	ID               int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre           string          `json:"nombre" gorm:"not null"`
	Descripcion      string          `json:"descripcion"`
	Unidad           string          `json:"unidad" gorm:"not null"`
	PrecioReferencia decimal.Decimal `json:"precio_referencia" gorm:"type:decimal(12,2);not null"`
	CreatedAt        time.Time       `json:"created_at" gorm:"autoCreateTime"`
	Precios          []PrecioRecurso `json:"precios,omitempty" gorm:"foreignKey:RecursoID;constraint:OnDelete:CASCADE"`
}

type EntradaMaterial struct {
	Nombre           string `json:"nombre"`
	Descripcion      string `json:"descripcion"`
	Unidad           string `json:"unidad"`
	PrecioReferencia string `json:"precio_referencia"`
}

var UnidadesPermitidas = map[string]bool{
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
	if !UnidadesPermitidas[m.Unidad] {
		return errors.New("unidad no permitida")
	}
	precio, err := decimal.NewFromString(m.PrecioReferencia)
	if err != nil || precio.LessThanOrEqual(decimal.Zero) {
		return errors.New("precio_referencia debe ser mayor a 0")
	}
	return nil
}

type ManoObra struct {
	ID              int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Descripcion     string          `json:"descripcion" gorm:"not null"`
	Categoria       string          `json:"categoria" gorm:"not null"`
	Unidad          string          `json:"unidad" gorm:"not null"`
	CostoReferencia decimal.Decimal `json:"costo_referencia" gorm:"type:decimal(12,2);not null"`
	CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime"`
	Precios         []PrecioRecurso `json:"precios,omitempty" gorm:"foreignKey:RecursoID;constraint:OnDelete:CASCADE"`
}

type EntradaManoObra struct {
	Descripcion     string          `json:"descripcion"`
	Categoria       string          `json:"categoria"`
	Unidad          string          `json:"unidad"`
	CostoReferencia decimal.Decimal `json:"costo_referencia"`
}

var CategoriasPermitidas = map[string]bool{
	"oficial": true, "ayudante": true, "especialista": true,
}

var UnidadesManoObra = map[string]bool{
	"hora": true, "día": true, "jornal": true,
}

func (m EntradaManoObra) ValidarManoObra() error {
	if m.Descripcion == "" {
		return errors.New("campo 'descripcion' requerido")
	}
	if !CategoriasPermitidas[m.Categoria] {
		return errors.New("categoria debe ser: oficial, ayudante, especialista")
	}
	if !UnidadesManoObra[m.Unidad] {
		return errors.New("unidad debe ser: hora, día, jornal")
	}
	if m.CostoReferencia.LessThanOrEqual(decimal.Zero) {
		return errors.New("costo_referencia debe ser mayor a 0")
	}
	return nil
}

type Equipo struct {
	ID         int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre     string          `json:"nombre" gorm:"not null"`
	Tipo       string          `json:"tipo" gorm:"not null"`
	Unidad     string          `json:"unidad" gorm:"not null"`
	CostoHora  decimal.Decimal `json:"costo_hora" gorm:"type:decimal(12,2);not null"`
	Disponible bool            `json:"disponible" gorm:"default:true"`
	CreatedAt  time.Time       `json:"created_at" gorm:"autoCreateTime"`
	Precios    []PrecioRecurso `json:"precios,omitempty" gorm:"foreignKey:RecursoID;constraint:OnDelete:CASCADE"`
}

type EntradaEquipo struct {
	Nombre     string          `json:"nombre"`
	Tipo       string          `json:"tipo"`
	Unidad     string          `json:"unidad"`
	CostoHora  decimal.Decimal `json:"costo_hora"`
	Disponible bool            `json:"disponible"`
}

var TiposEquipo = map[string]bool{
	"pesado": true, "liviano": true,
}

func (e EntradaEquipo) ValidarEquipo() error {
	if e.Nombre == "" {
		return errors.New("campo 'nombre' requerido")
	}
	if !TiposEquipo[e.Tipo] {
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

type PrecioRecurso struct {
	ID            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	RecursoTipo   string          `json:"recurso_tipo" gorm:"not null;index"`
	RecursoID     int             `json:"recurso_id" gorm:"not null;index"`
	Precio        decimal.Decimal `json:"precio" gorm:"type:decimal(12,2);not null"`
	FechaVigencia time.Time       `json:"fecha_vigencia" gorm:"not null"`
	CreatedAt     time.Time       `json:"created_at" gorm:"autoCreateTime"`
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
