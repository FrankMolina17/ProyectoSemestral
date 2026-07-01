package models

import "time"

type Obra struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Ubicacion   string    `json:"ubicacion"`
	Estado      string    `json:"estado"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}
