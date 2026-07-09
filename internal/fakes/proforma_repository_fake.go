package fakes

import (
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"
)

// ProformaRepositoryFake guarda en memoria para tests de handler (no usa la base real).
type ProformaRepositoryFake struct {
	mu            sync.Mutex
	proformas     map[int]models.Proforma
	items         map[int]models.ProformaItem
	notas         map[int]models.NotaProforma
	clientes      map[int]models.Cliente
	nextIDProf    int
	nextIDItem    int
	nextIDNota    int
	nextIDCliente int
}

func NuevoProformaRepositoryFake() *ProformaRepositoryFake {
	return &ProformaRepositoryFake{
		proformas:     make(map[int]models.Proforma),
		items:         make(map[int]models.ProformaItem),
		notas:         make(map[int]models.NotaProforma),
		clientes:      make(map[int]models.Cliente),
		nextIDProf:    1,
		nextIDItem:    1,
		nextIDNota:    1,
		nextIDCliente: 1,
	}
}

func (f *ProformaRepositoryFake) CrearProforma(p models.Proforma) (models.Proforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	p.ID = f.nextIDProf
	p.Estado = "borrador"
	p.CreadoEn = time.Now()
	p.Subtotal = 0
	p.Total = 0
	f.proformas[p.ID] = p
	f.nextIDProf++
	return p, nil
}

func (f *ProformaRepositoryFake) ObtenerPorID(id int) (models.Proforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	p, ok := f.proformas[id]
	if !ok {
		return models.Proforma{}, repository.ErrProformaNoEncontrada
	}
	return p, nil
}

func (f *ProformaRepositoryFake) ObtenerTodos() ([]models.Proforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	lista := make([]models.Proforma, 0, len(f.proformas))
	for _, p := range f.proformas {
		lista = append(lista, p)
	}
	return lista, nil
}

func (f *ProformaRepositoryFake) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	p, ok := f.proformas[id]
	if !ok {
		return models.Proforma{}, repository.ErrProformaNoEncontrada
	}
	p.Nombre = datos.Nombre
	p.PctGanancia = datos.PctGanancia
	p.PctImprevisto = datos.PctImprevisto
	p.ClienteID = datos.ClienteID
	f.proformas[id] = p
	return p, nil
}

func (f *ProformaRepositoryFake) EliminarProforma(id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.proformas[id]; !ok {
		return repository.ErrProformaNoEncontrada
	}
	delete(f.proformas, id)
	return nil
}

func (f *ProformaRepositoryFake) AprobarProforma(id int) (models.Proforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	p, ok := f.proformas[id]
	if !ok {
		return models.Proforma{}, repository.ErrProformaNoEncontrada
	}
	p.Estado = "aprobada"
	f.proformas[id] = p
	return p, nil
}

func (f *ProformaRepositoryFake) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.proformas[item.ProformaID]; !ok {
		return models.ProformaItem{}, repository.ErrProformaNoEncontrada
	}
	item.ID = f.nextIDItem
	item.Subtotal = item.Cantidad * item.PrecioPromedio
	f.items[item.ID] = item
	f.nextIDItem++
	f.recalcularTotales(item.ProformaID)
	return item, nil
}

func (f *ProformaRepositoryFake) ObtenerItems(proformaID int) ([]models.ProformaItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	lista := make([]models.ProformaItem, 0)
	for _, item := range f.items {
		if item.ProformaID == proformaID {
			lista = append(lista, item)
		}
	}
	return lista, nil
}

func (f *ProformaRepositoryFake) recalcularTotales(proformaID int) {
	p := f.proformas[proformaID]
	subtotal := 0.0
	for _, item := range f.items {
		if item.ProformaID == proformaID {
			subtotal += item.Subtotal
		}
	}
	p.Subtotal = subtotal
	p.Total = subtotal + (subtotal * p.PctGanancia) + (subtotal * p.PctImprevisto)
	f.proformas[proformaID] = p
}

func (f *ProformaRepositoryFake) AgregarNota(nota models.NotaProforma) (models.NotaProforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.proformas[nota.ProformaID]; !ok {
		return models.NotaProforma{}, repository.ErrProformaNoEncontrada
	}
	nota.ID = f.nextIDNota
	nota.CreadoEn = time.Now()
	f.notas[nota.ID] = nota
	f.nextIDNota++
	return nota, nil
}

func (f *ProformaRepositoryFake) ObtenerNotas(proformaID int) ([]models.NotaProforma, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	lista := make([]models.NotaProforma, 0)
	for _, n := range f.notas {
		if n.ProformaID == proformaID {
			lista = append(lista, n)
		}
	}
	return lista, nil
}

func (f *ProformaRepositoryFake) CrearCliente(c models.Cliente) (models.Cliente, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	c.ID = f.nextIDCliente
	f.clientes[c.ID] = c
	f.nextIDCliente++
	return c, nil
}

func (f *ProformaRepositoryFake) ObtenerClientes() ([]models.Cliente, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	lista := make([]models.Cliente, 0, len(f.clientes))
	for _, c := range f.clientes {
		lista = append(lista, c)
	}
	return lista, nil
}

func (f *ProformaRepositoryFake) ObtenerClientePorID(id int) (models.Cliente, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	c, ok := f.clientes[id]
	if !ok {
		return models.Cliente{}, repository.ErrClienteNoEncontrado
	}
	return c, nil
}

func (f *ProformaRepositoryFake) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	c, ok := f.clientes[id]
	if !ok {
		return models.Cliente{}, repository.ErrClienteNoEncontrado
	}
	c.Nombre = datos.Nombre
	c.Email = datos.Email
	c.Telefono = datos.Telefono
	c.Ruc = datos.Ruc
	f.clientes[id] = c
	return c, nil
}

func (f *ProformaRepositoryFake) EliminarCliente(id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.clientes[id]; !ok {
		return repository.ErrClienteNoEncontrado
	}
	delete(f.clientes, id)
	return nil
}
