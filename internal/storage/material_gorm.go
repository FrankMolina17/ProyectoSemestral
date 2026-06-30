package storage

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"gorm.io/gorm"
)

type MaterialGORM struct {
	db *gorm.DB
}

func NewMaterialGORM(db *gorm.DB) *MaterialGORM {
	return &MaterialGORM{db: db}
}

func (r *MaterialGORM) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	mat := models.Material{
		Nombre:           in.Nombre,
		Descripcion:      in.Descripcion,
		Unidad:           in.Unidad,
		PrecioReferencia: in.PrecioReferencia,
	}
	if err := r.db.Create(&mat).Error; err != nil {
		return nil, err
	}
	return &mat, nil
}

func (r *MaterialGORM) ObtenerMateriales(id int) (*models.Material, bool) {
	var mat models.Material
	if err := r.db.First(&mat, id).Error; err != nil {
		return nil, false
	}
	return &mat, true
}

func (r *MaterialGORM) ListarMateriales() []*models.Material {
	var mats []models.Material
	r.db.Find(&mats)
	out := make([]*models.Material, len(mats))
	for i := range mats {
		out[i] = &mats[i]
	}
	return out
}

func (r *MaterialGORM) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	var mat models.Material
	if err := r.db.First(&mat, id).Error; err != nil {
		return nil, false
	}
	mat.Nombre = in.Nombre
	mat.Descripcion = in.Descripcion
	mat.Unidad = in.Unidad
	mat.PrecioReferencia = in.PrecioReferencia
	if err := r.db.Save(&mat).Error; err != nil {
		return nil, false
	}
	return &mat, true
}

func (r *MaterialGORM) EliminarMateriales(id int) bool {
	res := r.db.Delete(&models.Material{}, id)
	if res.Error != nil || res.RowsAffected == 0 {
		return false
	}
	return true
}
