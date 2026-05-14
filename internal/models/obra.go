package models

type Obra struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Ubicacion   string `json:"ubicacion"`
	Estado      string `json:"estado"`
}

type Incidencia struct {
	ID          int    `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Tipo        string `json:"tipo"`
	Estado      string `json:"estado"`
}

// ejemplo para commit