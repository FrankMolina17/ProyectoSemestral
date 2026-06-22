package models

import "time"

type Proforma struct {
    ID             int       `json:"id"`
    ObraID         int       `json:"obra_id"`
    ClienteID      int       `json:"cliente_id"`
    Nombre         string    `json:"nombre"`
    Estado         string    `json:"estado"`
    PctGanancia    float64   `json:"pct_ganancia"`
    PctImprevisto  float64   `json:"pct_imprevisto"`
    Subtotal       float64   `json:"subtotal"`
    Total          float64   `json:"total"`
    CreadoEn       time.Time `json:"creado_en"`
}

type ProformaItem struct {
    ID             int     `json:"id"`
    ProformaID     int     `json:"proforma_id"`
    TipoRecurso    string  `json:"tipo_recurso"`
    RecursoID      int     `json:"recurso_id"`
    Descripcion    string  `json:"descripcion"`
    Cantidad       float64 `json:"cantidad"`
    PrecioPromedio float64 `json:"precio_promedio"`
    Subtotal       float64 `json:"subtotal"`
}

type Cliente struct {
    ID       int    `json:"id"`
    Nombre   string `json:"nombre"`
    Email    string `json:"email"`
    Telefono string `json:"telefono"`
    Ruc      string `json:"ruc"`
}

type NotaProforma struct {
    ID         int       `json:"id"`
    ProformaID int       `json:"proforma_id"`
    Contenido  string    `json:"contenido"`
    CreadoEn   time.Time `json:"creado_en"`
}