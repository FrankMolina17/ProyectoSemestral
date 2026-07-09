package models

import "time"

type Material struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Unidad      string  `json:"unidad"`
	Cantidad    float64 `json:"cantidad"`
	CostoUnit   float64 `json:"costo_unitario"`
	Proveedor   string  `json:"proveedor,omitempty"`
}

type ManoObra struct {
	ID            int     `json:"id"`
	Nombre        string  `json:"nombre"`
	Tipo          string  `json:"tipo"`
	CostoPorHora  float64 `json:"costo_por_hora"`
}

type Equipo struct {
	ID              int     `json:"id"`
	Nombre          string  `json:"nombre"`
	Modelo          string  `json:"modelo,omitempty"`
	CostoPorHora    float64 `json:"costo_por_hora"`
}

type Precio struct {
	ID        int       `json:"id"`
	Tipo      string    `json:"tipo"`
	RecursoID int       `json:"recurso_id"`
	Monto     float64   `json:"monto"`
	Vigente   bool      `json:"vigente"`
	CreadoEn  time.Time `json:"creado_en"`
}

type UsuarioCatalogo struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Nombre   string `json:"nombre,omitempty"`
}
