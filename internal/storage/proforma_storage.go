package storage

import (
    "errors"
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
    "sync"
    "time"
)

// Storage guarda todo en memoria
type ProformaStorage struct {
    mu          sync.Mutex
    proformas   map[int]models.Proforma
    items       map[int]models.ProformaItem
    nextIDProf  int
    nextIDItem  int
}

// New crea un storage vacío listo para usar
func NuevoStorage() *ProformaStorage {
    return &ProformaStorage{
        proformas:  make(map[int]models.Proforma),
        items:      make(map[int]models.ProformaItem),
        nextIDProf: 1,
        nextIDItem: 1,
    }
}


func (s *ProformaStorage) CrearProforma(p models.Proforma) models.Proforma {
    s.mu.Lock()
    defer s.mu.Unlock()

    p.ID = s.nextIDProf
    p.Estado = "borrador"
    p.CreadoEn = time.Now()
    p.Subtotal = 0
    p.Total = 0

    s.proformas[p.ID] = p
    s.nextIDProf++
    return p
}

func (s *ProformaStorage) ObtenerPorID(id int) (models.Proforma, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    p, ok := s.proformas[id]
    if !ok {
        return models.Proforma{}, errors.New("proforma no encontrada")
    }
    return p, nil
}

func (s *ProformaStorage) ObtenerTodos() []models.Proforma {
    s.mu.Lock()
    defer s.mu.Unlock()

    lista := make([]models.Proforma, 0, len(s.proformas))
    for _, p := range s.proformas {
        lista = append(lista, p)
    }
    return lista
}

func (s *ProformaStorage) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    p, ok := s.proformas[id]
    if !ok {
        return models.Proforma{}, errors.New("proforma no encontrada")
    }

    // Solo actualizas los campos editables
    p.Nombre = datos.Nombre
    p.PctGanancia = datos.PctGanancia
    p.PctImprevisto = datos.PctImprevisto

    s.proformas[id] = p
    return p, nil
}

func (s *ProformaStorage) EliminarProforma(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    _, ok := s.proformas[id]
    if !ok {
        return errors.New("proforma no encontrada")
    }

    delete(s.proformas, id)
    return nil
}

func (s *ProformaStorage) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Verificar que la proforma existe
    _, ok := s.proformas[item.ProformaID]
    if !ok {
        return models.ProformaItem{}, errors.New("proforma no encontrada")
    }

    // Calcular subtotal del ítem
    item.ID = s.nextIDItem
    item.Subtotal = item.Cantidad * item.PrecioPromedio

    s.items[item.ID] = item
    s.nextIDItem++

    // Recalcular totales de la proforma
    s.recalcularTotales(item.ProformaID)

    return item, nil
}

func (s *ProformaStorage) ObtenerItems(proformaID int) []models.ProformaItem {
    s.mu.Lock()
    defer s.mu.Unlock()

    lista := make([]models.ProformaItem, 0)
    for _, item := range s.items {
        if item.ProformaID == proformaID {
            lista = append(lista, item)
        }
    }
    return lista
}

func (s *ProformaStorage) recalcularTotales(proformaID int) {
    p := s.proformas[proformaID]
    subtotal := 0.0

    for _, item := range s.items {
        if item.ProformaID == proformaID {
            subtotal += item.Subtotal
        }
    }

    p.Subtotal = subtotal
    p.Total = subtotal + (subtotal * p.PctGanancia) + (subtotal * p.PctImprevisto)
    s.proformas[proformaID] = p
}