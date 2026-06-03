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

func CrearIncidencia(incidencia models.Incidencia) (models.Incidencia, error) {
	if incidencia.Titulo == "" {
		return incidencia, errors.New("el titulo es obligatorio")
	}
	if incidencia.EntidadTipo == "" || incidencia.EntidadID <= 0 {
		return incidencia, errors.New("entidad_tipo y entidad_id son obligatorios")
	}

	// Asignar ID automático
	incidencia.ID = storage.IncidenciaIDCounter
	storage.IncidenciaIDCounter++

	// Valor por defecto de estado
	if incidencia.Estado == "" {
		incidencia.Estado = "abierta"
	}

	storage.Incidencias = append(storage.Incidencias, incidencia)
	return incidencia, nil
}

func ObtenerObras() []models.Obra {
	return storage.Obras
}

func ObtenerIncidencias() []models.Incidencia {
	return storage.Incidencias
}

func ObtenerIncidenciasPorEntidad(entidadTipo string, entidadID int) []models.Incidencia {
	var resultado []models.Incidencia
	for _, inc := range storage.Incidencias {
		if inc.EntidadTipo == entidadTipo && inc.EntidadID == entidadID {
			resultado = append(resultado, inc)
		}
	}
	return resultado
}

func ObtenerIncidenciaPorID(id int) (models.Incidencia, bool) {
	for _, inc := range storage.Incidencias {
		if inc.ID == id {
			return inc, true
		}
	}
	return models.Incidencia{}, false
}

func ActualizarIncidencia(id int, incidencia models.Incidencia) (models.Incidencia, error) {
	for i, inc := range storage.Incidencias {
		if inc.ID == id {
			// Actualizar campos
			if incidencia.Titulo != "" {
				storage.Incidencias[i].Titulo = incidencia.Titulo
			}
			if incidencia.Descripcion != "" {
				storage.Incidencias[i].Descripcion = incidencia.Descripcion
			}
			if incidencia.Estado != "" {
				storage.Incidencias[i].Estado = incidencia.Estado
			}
			return storage.Incidencias[i], nil
		}
	}	
	return models.Incidencia{}, errors.New("incidencia no encontrada")
}


func eliminarIncidencia(id int) error {
	for i, inc := range storage.Incidencias {
		if inc.ID == id {			
			// Eliminar incidencia del slice
			storage.Incidencias = append(storage.Incidencias[:i], storage.Incidencias[i+1:]...)
			return nil
		}	
	}
	return errors.New("incidencia no encontrada")
}
