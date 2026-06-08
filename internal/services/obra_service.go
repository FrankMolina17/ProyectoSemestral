package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"errors"
)

func CrearObra(obra models.Obra) (models.Obra, error) {
	if obra.Nombre == "" {
		return obra, errors.New("el nombre de la obra es obligatorio")
	}
	if obra.UserID <= 0 {
		return obra, errors.New("el user_id es obligatorio")
	}

	// Asignar ID automático
	obra.ID = storage.ObraIDCounter
	storage.ObraIDCounter++

	// Valor por defecto de estado
	if obra.Estado == "" {
		obra.Estado = "planificacion"
	}

	storage.Obras = append(storage.Obras, obra)
	return obra, nil
}

func ObtenerObras() []models.Obra {
	return storage.Obras
}

func ObtenerObra(id int) (models.Obra, bool) {
	for _, obra := range storage.Obras {
		if obra.ID == id {
			return obra, true
		}
	}

	return models.Obra{}, false
}

func ActualizarObra(id int, obra models.Obra) (models.Obra, bool) {
	for i, o := range storage.Obras {
		if o.ID == id {
			if obra.Nombre != "" {
				storage.Obras[i].Nombre = obra.Nombre
			}
			if obra.Descripcion != "" {
				storage.Obras[i].Descripcion = obra.Descripcion
			}
			if obra.Ubicacion != "" {
				storage.Obras[i].Ubicacion = obra.Ubicacion
			}
			return storage.Obras[i], true
		}
	}
	return models.Obra{}, false
}

func EliminarObra(id int) bool {
	for i, o := range storage.Obras {
		if o.ID == id {
			storage.Obras = append(storage.Obras[:i], storage.Obras[i+1:]...)
			return true
		}
	}
	return false
}
