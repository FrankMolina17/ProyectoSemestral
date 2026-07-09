package models

type Incidencia struct {
	ID            int     `json:"id" gorm:"primaryKey;autoIncrement"`
	EntidadTipo   string  `json:"entidad_tipo" gorm:"not null;index"` // "obra" o "proforma"
	EntidadID     int     `json:"entidad_id" gorm:"not null;index"`
	ResponsableID int     `json:"responsable_id" gorm:"index"`
	Titulo        string  `json:"titulo" gorm:"not null"`
	Descripcion   string  `json:"descripcion"`
	Tipo          string  `json:"tipo" gorm:"not null"`
	Prioridad     string  `json:"prioridad" gorm:"default:media"`
	Estado        string  `json:"estado" gorm:"default:abierta"`
	ImpactoCosto  float64 `json:"impacto_costo"`
	ImpactoTiempo int     `json:"impacto_tiempo"`
}
