package models

type Incidencia struct {
	ID            int    `json:"id" gorm:"primaryKey"`
	EntidadTipo   string `json:"entidad_tipo"` // "obra" o "proforma"
	EntidadID     int    `json:"entidad_id"`
	ResponsableID int    `json:"responsable_id"`
	Titulo        string `json:"titulo"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
}
