package models

import "time"

type Obra struct {
	ID                  int         `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre              string      `json:"nombre" gorm:"not null"`
	Descripcion         string      `json:"descripcion"`
	Ubicacion           string      `json:"ubicacion"`
	FechaInicio         time.Time   `json:"fecha_inicio"`
	FechaFin            time.Time   `json:"fecha_fin"`
	Estado              string      `json:"estado" gorm:"default:planificada"`
	PresupuestoEstimado float64     `json:"presupuesto_estimado"`
	UserID              int         `json:"user_id" gorm:"index"`
	ProformaID          *int        `json:"proforma_id" gorm:"index"`
	Incidencias         []Incidencia `json:"incidencias,omitempty" gorm:"foreignKey:EntidadID;constraint:OnDelete:CASCADE"`
}
