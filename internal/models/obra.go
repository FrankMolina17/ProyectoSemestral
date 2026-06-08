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
