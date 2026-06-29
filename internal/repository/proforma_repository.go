package repository

import (
	"errors"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ErrProformaNoEncontrada = errors.New("proforma no encontrada")
var ErrClienteNoEncontrado = errors.New("cliente no encontrado")

type ProformaRepository interface {
	CrearProforma(p models.Proforma) (models.Proforma, error)
	ObtenerPorID(id int) (models.Proforma, error)
	ObtenerTodos() ([]models.Proforma, error)
	ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error)
	EliminarProforma(id int) error
	AprobarProforma(id int) (models.Proforma, error)
	AgregarItem(item models.ProformaItem) (models.ProformaItem, error)
	ObtenerItems(proformaID int) ([]models.ProformaItem, error)
	AgregarNota(nota models.NotaProforma) (models.NotaProforma, error)
	ObtenerNotas(proformaID int) ([]models.NotaProforma, error)
	CrearCliente(c models.Cliente) (models.Cliente, error)
	ObtenerClientes() ([]models.Cliente, error)
	ObtenerClientePorID(id int) (models.Cliente, error)
	ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error)
	EliminarCliente(id int) error
}

type proformaGormRepo struct {
	db *gorm.DB
}

func NuevaConexion(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Proforma{},
		&models.ProformaItem{},
		&models.Cliente{},
		&models.NotaProforma{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

func NuevoProformaRepository(db *gorm.DB) ProformaRepository {
	return &proformaGormRepo{db: db}
}

func (r *proformaGormRepo) CrearProforma(p models.Proforma) (models.Proforma, error) {
	p.Estado = "borrador"
	p.CreadoEn = time.Now()
	p.Subtotal = 0
	p.Total = 0

	if err := r.db.Create(&p).Error; err != nil {
		return models.Proforma{}, err
	}
	return p, nil
}

func (r *proformaGormRepo) ObtenerPorID(id int) (models.Proforma, error) {
	var p models.Proforma
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Proforma{}, ErrProformaNoEncontrada
		}
		return models.Proforma{}, err
	}
	return p, nil
}

func (r *proformaGormRepo) ObtenerTodos() ([]models.Proforma, error) {
	var lista []models.Proforma
	if err := r.db.Order("id asc").Find(&lista).Error; err != nil {
		return nil, err
	}
	return lista, nil
}

func (r *proformaGormRepo) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
	p, err := r.ObtenerPorID(id)
	if err != nil {
		return models.Proforma{}, err
	}

	p.Nombre = datos.Nombre
	p.PctGanancia = datos.PctGanancia
	p.PctImprevisto = datos.PctImprevisto
	p.ClienteID = datos.ClienteID

	if err := r.db.Save(&p).Error; err != nil {
		return models.Proforma{}, err
	}
	return p, nil
}

func (r *proformaGormRepo) EliminarProforma(id int) error {
	result := r.db.Delete(&models.Proforma{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProformaNoEncontrada
	}
	return nil
} 

func (r *proformaGormRepo) AprobarProforma(id int) (models.Proforma, error) {
	p, err := r.ObtenerPorID(id)
	if err != nil {
		return models.Proforma{}, err
	}

	p.Estado = "aprobada"
	if err := r.db.Save(&p).Error; err != nil {
		return models.Proforma{}, err
	}
	return p, nil
}

func (r *proformaGormRepo) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
	if _, err := r.ObtenerPorID(item.ProformaID); err != nil {
		return models.ProformaItem{}, err
	}

	item.Subtotal = item.Cantidad * item.PrecioPromedio
	if err := r.db.Create(&item).Error; err != nil {
		return models.ProformaItem{}, err
	}

	if err := r.recalcularTotales(item.ProformaID); err != nil {
		return models.ProformaItem{}, err
	}
	return item, nil
}

func (r *proformaGormRepo) ObtenerItems(proformaID int) ([]models.ProformaItem, error) {
	var items []models.ProformaItem
	if err := r.db.Where("proforma_id = ?", proformaID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *proformaGormRepo) recalcularTotales(proformaID int) error {
	p, err := r.ObtenerPorID(proformaID)
	if err != nil {
		return err
	}

	items, err := r.ObtenerItems(proformaID)
	if err != nil {
		return err
	}

	subtotal := 0.0
	for _, item := range items {
		subtotal += item.Subtotal
	}

	p.Subtotal = subtotal
	p.Total = subtotal + (subtotal * p.PctGanancia) + (subtotal * p.PctImprevisto)
	return r.db.Save(&p).Error
}

func (r *proformaGormRepo) AgregarNota(nota models.NotaProforma) (models.NotaProforma, error) {
	if _, err := r.ObtenerPorID(nota.ProformaID); err != nil {
		return models.NotaProforma{}, err
	}

	nota.CreadoEn = time.Now()
	if err := r.db.Create(&nota).Error; err != nil {
		return models.NotaProforma{}, err
	}
	return nota, nil
}

func (r *proformaGormRepo) ObtenerNotas(proformaID int) ([]models.NotaProforma, error) {
	var notas []models.NotaProforma
	if err := r.db.Where("proforma_id = ?", proformaID).Find(&notas).Error; err != nil {
		return nil, err
	}
	return notas, nil
}

func (r *proformaGormRepo) CrearCliente(c models.Cliente) (models.Cliente, error) {
	if err := r.db.Create(&c).Error; err != nil {
		return models.Cliente{}, err
	}
	return c, nil
}

func (r *proformaGormRepo) ObtenerClientes() ([]models.Cliente, error) {
	var clientes []models.Cliente
	if err := r.db.Order("id asc").Find(&clientes).Error; err != nil {
		return nil, err
	}
	return clientes, nil
}

func (r *proformaGormRepo) ObtenerClientePorID(id int) (models.Cliente, error) {
	var c models.Cliente
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cliente{}, ErrClienteNoEncontrado
		}
		return models.Cliente{}, err
	}
	return c, nil
}

func (r *proformaGormRepo) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error) {
	c, err := r.ObtenerClientePorID(id)
	if err != nil {
		return models.Cliente{}, err
	}

	c.Nombre = datos.Nombre
	c.Email = datos.Email
	c.Telefono = datos.Telefono
	c.Ruc = datos.Ruc

	if err := r.db.Save(&c).Error; err != nil {
		return models.Cliente{}, err
	}
	return c, nil
}

func (r *proformaGormRepo) EliminarCliente(id int) error {
	result := r.db.Delete(&models.Cliente{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClienteNoEncontrado
	}
	return nil
}
