package models

// Proforma representa un presupuesto estimado para una obra.
type Proforma struct {
	ID          int            `json:"id"`
	ObraID      int            `json:"obra_id"`
	Nombre      string         `json:"nombre"`
	Descripcion string         `json:"descripcion"`
	Estado      string         `json:"estado"`
	Items       []ItemProforma `json:"items"`
	Subtotal    float64        `json:"subtotal"`
	Impuestos   float64        `json:"impuestos"`
	Total       float64        `json:"total"`
}

// ItemProforma representa un recurso dentro de una proforma con su costo.
type ItemProforma struct {
	ID          int     `json:"id"`
	RecursoID   int     `json:"recurso_id"`
	Tipo        string  `json:"tipo"`
	Descripcion string  `json:"descripcion"`
	Cantidad    float64 `json:"cantidad"`
	Unidad      string  `json:"unidad"`
	PrecioUnit  float64 `json:"precio_unitario"`
	Subtotal    float64 `json:"subtotal"`
}
