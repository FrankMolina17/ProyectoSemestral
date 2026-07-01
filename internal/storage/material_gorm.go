package storage

import (
	"errors"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MaterialGORM struct {
	db *gorm.DB
}

func NewMaterialGORM(db *gorm.DB) *MaterialGORM {
	return &MaterialGORM{db: db}
}

func precioADecimal(precioStr string) (decimal.Decimal, error) {
	if precioStr == "" {
		return decimal.Zero, errors.New("precio vacio")
	}
	precio, err := decimal.NewFromString(precioStr)
	if err != nil {
		return decimal.Zero, err
	}
	return precio, nil
}

func (r *MaterialGORM) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	precio, err := precioADecimal(in.PrecioReferencia)
	if err != nil {
		return nil, err
	}
	mat := models.Material{
		Nombre:           in.Nombre,
		Descripcion:      in.Descripcion,
		Unidad:           in.Unidad,
		PrecioReferencia: precio,
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
	precio, err := precioADecimal(in.PrecioReferencia)
	if err != nil {
		return nil, false
	}
	mat.Nombre = in.Nombre
	mat.Descripcion = in.Descripcion
	mat.Unidad = in.Unidad
	mat.PrecioReferencia = precio
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
