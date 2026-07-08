package models

import "time"

type Obra struct {
	ID                  int       `json:"id" gorm:"primaryKey"`
	Nombre              string    `json:"nombre" gorm:"not null"`
	Descripcion         string    `json:"descripcion"`
	Ubicacion           string    `json:"ubicacion"`
	Estado              string    `json:"estado"`
	UserID              int       `json:"user_id" gorm:"not null"`
	FechaInicio         time.Time `json:"fecha_inicio"`
	FechaFinEstimada    time.Time `json:"fecha_fin_estimada"`
	FechaFinReal        time.Time `json:"fecha_fin_real,omitempty"`
	PresupuestoEstimado float64   `json:"presupuesto_estimado"`
	PresupuestoReal     float64   `json:"presupuesto_real,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
