package models

import "time"

type Proforma struct {
    ID             int       `json:"id"`
    ObraID         int       `json:"obra_id"`
    Nombre         string    `json:"nombre"`
    Estado         string    `json:"estado"`          // "borrador" o "aprobada"
    PctGanancia    float64   `json:"pct_ganancia"`    // ej: 0.10 = 10%
    PctImprevisto  float64   `json:"pct_imprevisto"`  // ej: 0.05 = 5%
    Subtotal       float64   `json:"subtotal"`
    Total          float64   `json:"total"`
    CreadoEn       time.Time `json:"creado_en"`
}

type ProformaItem struct {
    ID             int     `json:"id"`
    ProformaID     int     `json:"proforma_id"`
    TipoRecurso    string  `json:"tipo_recurso"`    // "material", "mano_obra", "equipo"
    RecursoID      int     `json:"recurso_id"`
    Descripcion    string  `json:"descripcion"`
    Cantidad       float64 `json:"cantidad"`
    PrecioPromedio float64 `json:"precio_promedio"` // calculado por el sistema
    Subtotal       float64 `json:"subtotal"`        // cantidad × precio_promedio
}