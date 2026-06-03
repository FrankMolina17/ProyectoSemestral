package models

import "time"

type Obra struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	Nombre              string    `json:"nombre"`
	Descripcion         string    `json:"descripcion"`
	Ubicacion           string    `json:"ubicacion"`
	FechaInicio         time.Time `json:"fecha_inicio"`
	FechaFin            time.Time `json:"fecha_fin"`
	Estado              string    `json:"estado"`
	PresupuestoEstimado float64   `json:"presupuesto_estimado"`
	UserID              int       `json:"user_id"`
	ProformaID          *int      `json:"proforma_id"`
}

type Incidencia struct {
	ID            int     `json:"id" gorm:"primaryKey"`
	EntidadTipo   string  `json:"entidad_tipo"` // "obra" o "proforma"
	EntidadID     int     `json:"entidad_id"`
	ResponsableID int     `json:"responsable_id"`
	Titulo        string  `json:"titulo"`
	Descripcion   string  `json:"descripcion"`
	Tipo          string  `json:"tipo"`
	Prioridad     string  `json:"prioridad"`
	Estado        string  `json:"estado"`
	ImpactoCosto  float64 `json:"impacto_costo"`
	ImpactoTiempo int     `json:"impacto_tiempo"`
}
