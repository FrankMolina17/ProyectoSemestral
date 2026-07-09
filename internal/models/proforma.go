package models

import "time"

type Proforma struct {
	ID            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	ObraID        int             `json:"obra_id" gorm:"not null"`
	ClienteID     int             `json:"cliente_id" gorm:"index"`
	Nombre        string          `json:"nombre" gorm:"not null"`
	Estado        string          `json:"estado" gorm:"default:borrador"`
	PctGanancia   float64         `json:"pct_ganancia"`
	PctImprevisto float64         `json:"pct_imprevisto"`
	Subtotal      float64         `json:"subtotal"`
	Total         float64         `json:"total"`
	CreadoEn      time.Time       `json:"creado_en" gorm:"autoCreateTime"`
	Cliente       *Cliente        `json:"cliente,omitempty" gorm:"foreignKey:ClienteID"`
	Items         []ProformaItem  `json:"items,omitempty" gorm:"foreignKey:ProformaID"`
	Notas         []NotaProforma  `json:"notas,omitempty" gorm:"foreignKey:ProformaID"`
	Incidencias   []Incidencia    `json:"incidencias,omitempty" gorm:"foreignKey:EntidadID;constraint:OnDelete:CASCADE"`
}

type ProformaItem struct {
	ID             int     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProformaID     int     `json:"proforma_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	TipoRecurso    string  `json:"tipo_recurso"`
	RecursoID      int     `json:"recurso_id"`
	Descripcion    string  `json:"descripcion" gorm:"not null"`
	Cantidad       float64 `json:"cantidad" gorm:"not null"`
	PrecioPromedio float64 `json:"precio_promedio" gorm:"not null"`
	Subtotal       float64 `json:"subtotal"`
}

type Cliente struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nombre   string `json:"nombre" gorm:"not null"`
	Email    string `json:"email"`
	Telefono string `json:"telefono"`
	Ruc      string `json:"ruc" gorm:"not null;uniqueIndex"`
}

type NotaProforma struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ProformaID int       `json:"proforma_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Contenido  string    `json:"contenido" gorm:"not null"`
	CreadoEn   time.Time `json:"creado_en" gorm:"autoCreateTime"`
}
