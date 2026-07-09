package models

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
