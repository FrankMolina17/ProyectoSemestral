package services

import (
	"errors"
	"testing"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────
// Material
// ─────────────────────────────────────────────

type materialRepoFull struct {
	materiales []*models.Material
	nextID     int
}

func (m *materialRepoFull) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	m.nextID++
	precio, _ := decimal.NewFromString(in.PrecioReferencia)
	mat := &models.Material{
		ID:               m.nextID,
		Nombre:           in.Nombre,
		Descripcion:      in.Descripcion,
		Unidad:           in.Unidad,
		PrecioReferencia: precio,
	}
	m.materiales = append(m.materiales, mat)
	return mat, nil
}

func (m *materialRepoFull) ObtenerMateriales(id int) (*models.Material, bool) {
	for _, mat := range m.materiales {
		if mat.ID == id {
			return mat, true
		}
	}
	return nil, false
}

func (m *materialRepoFull) ListarMateriales() []*models.Material {
	if m.materiales == nil {
		return []*models.Material{}
	}
	return m.materiales
}

func (m *materialRepoFull) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	for i, mat := range m.materiales {
		if mat.ID == id {
			m.materiales[i].Nombre = in.Nombre
			m.materiales[i].Descripcion = in.Descripcion
			m.materiales[i].Unidad = in.Unidad
			m.materiales[i].PrecioReferencia, _ = decimal.NewFromString(in.PrecioReferencia)
			return m.materiales[i], true
		}
	}
	return nil, false
}

func (m *materialRepoFull) EliminarMateriales(id int) bool {
	for i, mat := range m.materiales {
		if mat.ID == id {
			m.materiales = append(m.materiales[:i], m.materiales[i+1:]...)
			return true
		}
	}
	return false
}

func TestMaterialService_CRUD(t *testing.T) {
	svc := NewMaterialService(&materialRepoFull{})

	t.Run("Listado vacio", func(t *testing.T) {
		res := svc.Listado()
		assert.Empty(t, res)
	})

	t.Run("Obtener inexistente -> error", func(t *testing.T) {
		_, err := svc.ObtenerM(999)
		assert.Error(t, err)
	})

	t.Run("Crear y luego obtener", func(t *testing.T) {
		in := models.EntradaMaterial{
			Nombre:           "Ladrillo",
			Descripcion:      "Rojo 10x20",
			Unidad:           "unidad",
			PrecioReferencia: "0.50",
		}
		creado, err := svc.CrearM(in)
		require.NoError(t, err)
		assert.Equal(t, "Ladrillo", creado.Nombre)

		obtenido, err := svc.ObtenerM(creado.ID)
		require.NoError(t, err)
		assert.Equal(t, creado.ID, obtenido.ID)
	})

	t.Run("Actualizar", func(t *testing.T) {
		in := models.EntradaMaterial{
			Nombre:           "Ladrillo Hueco",
			Descripcion:      "Hueco 10x20",
			Unidad:           "unidad",
			PrecioReferencia: "0.60",
		}
		actualizado, err := svc.ActualizarM(1, in)
		require.NoError(t, err)
		assert.Equal(t, "Ladrillo Hueco", actualizado.Nombre)
	})

	t.Run("Eliminar", func(t *testing.T) {
		err := svc.EliminarM(1)
		assert.NoError(t, err)

		_, err = svc.ObtenerM(1)
		assert.Error(t, err)
	})

	t.Run("Actualizar inexistente -> error", func(t *testing.T) {
		_, err := svc.ActualizarM(999, models.EntradaMaterial{})
		assert.Error(t, err)
	})

	t.Run("Eliminar inexistente -> error", func(t *testing.T) {
		err := svc.EliminarM(999)
		assert.Error(t, err)
	})
}

// ─────────────────────────────────────────────
// ManoObra
// ─────────────────────────────────────────────

type manoObraRepoFull struct {
	items  []*models.ManoObra
	nextID int
}

func (m *manoObraRepoFull) CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error) {
	m.nextID++
	item := &models.ManoObra{
		ID:              m.nextID,
		Descripcion:     in.Descripcion,
		Categoria:       in.Categoria,
		Unidad:          in.Unidad,
		CostoReferencia: in.CostoReferencia,
	}
	m.items = append(m.items, item)
	return item, nil
}

func (m *manoObraRepoFull) ObtenerManoObra(id int) (*models.ManoObra, bool) {
	for _, item := range m.items {
		if item.ID == id {
			return item, true
		}
	}
	return nil, false
}

func (m *manoObraRepoFull) ListarManoObra() []*models.ManoObra {
	if m.items == nil {
		return []*models.ManoObra{}
	}
	return m.items
}

func (m *manoObraRepoFull) ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool) {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i].Descripcion = in.Descripcion
			m.items[i].Categoria = in.Categoria
			m.items[i].Unidad = in.Unidad
			m.items[i].CostoReferencia = in.CostoReferencia
			return m.items[i], true
		}
	}
	return nil, false
}

func (m *manoObraRepoFull) EliminarManoObra(id int) bool {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return true
		}
	}
	return false
}

func TestManoObraService_CRUD(t *testing.T) {
	svc := NewManoObraService(&manoObraRepoFull{})

	t.Run("Listado vacio", func(t *testing.T) {
		assert.Empty(t, svc.ListadoMa())
	})

	t.Run("Obtener inexistente -> error", func(t *testing.T) {
		_, err := svc.ObtenerMa(1)
		assert.Error(t, err)
	})

	t.Run("Crear", func(t *testing.T) {
		in := models.EntradaManoObra{
			Descripcion:     "Albañil",
			Categoria:       "oficial",
			Unidad:          "hora",
			CostoReferencia: decimal.NewFromFloat(15.00),
		}
		creado, err := svc.CrearMa(in)
		require.NoError(t, err)
		assert.Equal(t, "Albañil", creado.Descripcion)
	})

	t.Run("Actualizar", func(t *testing.T) {
		in := models.EntradaManoObra{
			Descripcion:     "Albañil Senior",
			Categoria:       "oficial",
			Unidad:          "hora",
			CostoReferencia: decimal.NewFromFloat(20.00),
		}
		actualizado, err := svc.ActualizarMa(1, in)
		require.NoError(t, err)
		assert.Equal(t, "Albañil Senior", actualizado.Descripcion)
	})

	t.Run("Eliminar", func(t *testing.T) {
		assert.True(t, svc.EliminarMa(1))
		assert.False(t, svc.EliminarMa(1))
	})
}

// ─────────────────────────────────────────────
// Equipo
// ─────────────────────────────────────────────

type equipoRepoFull struct {
	items  []*models.Equipo
	nextID int
}

func (e *equipoRepoFull) CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error) {
	e.nextID++
	item := &models.Equipo{
		ID:        e.nextID,
		Nombre:    in.Nombre,
		Tipo:      in.Tipo,
		Unidad:    in.Unidad,
		CostoHora: in.CostoHora,
	}
	e.items = append(e.items, item)
	return item, nil
}

func (e *equipoRepoFull) ObtenerEquipo(id int) (*models.Equipo, error) {
	for _, item := range e.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, errors.New("not found")
}

func (e *equipoRepoFull) ListarEquipos() []*models.Equipo {
	if e.items == nil {
		return []*models.Equipo{}
	}
	return e.items
}

func (e *equipoRepoFull) ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error) {
	for i, item := range e.items {
		if item.ID == id {
			e.items[i].Nombre = in.Nombre
			e.items[i].Tipo = in.Tipo
			e.items[i].Unidad = in.Unidad
			e.items[i].CostoHora = in.CostoHora
			return e.items[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (e *equipoRepoFull) EliminarEquipo(id int) error {
	for i, item := range e.items {
		if item.ID == id {
			e.items = append(e.items[:i], e.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestEquipoService_CRUD(t *testing.T) {
	svc := NewEquipoService(&equipoRepoFull{})

	t.Run("Listado vacio", func(t *testing.T) {
		assert.Empty(t, svc.ListadoE())
	})

	t.Run("Obtener inexistente -> error", func(t *testing.T) {
		_, err := svc.ObtenerE(1)
		assert.Error(t, err)
	})

	t.Run("Crear", func(t *testing.T) {
		in := models.EntradaEquipo{
			Nombre:    "Excavadora",
			Tipo:      "pesado",
			Unidad:    "hora",
			CostoHora: decimal.NewFromFloat(85.00),
		}
		creado, err := svc.CrearE(in)
		require.NoError(t, err)
		assert.Equal(t, "Excavadora", creado.Nombre)
	})

	t.Run("Actualizar", func(t *testing.T) {
		in := models.EntradaEquipo{
			Nombre:    "Excavadora Grande",
			Tipo:      "pesado",
			Unidad:    "hora",
			CostoHora: decimal.NewFromFloat(95.00),
		}
		actualizado, err := svc.ActualizarE(1, in)
		require.NoError(t, err)
		assert.Equal(t, "Excavadora Grande", actualizado.Nombre)
	})

	t.Run("Eliminar", func(t *testing.T) {
		err := svc.EliminarE(1)
		assert.NoError(t, err)

		err = svc.EliminarE(1)
		assert.Error(t, err)
	})
}

// ─────────────────────────────────────────────
// Precio
// ─────────────────────────────────────────────

type precioRepoFull struct {
	items  []*models.PrecioRecurso
	nextID int
}

func (p *precioRepoFull) CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	p.nextID++
	item := &models.PrecioRecurso{
		ID:            p.nextID,
		RecursoTipo:   in.RecursoTipo,
		RecursoID:     in.RecursoID,
		Precio:        in.Precio,
		FechaVigencia: in.FechaVigencia,
	}
	p.items = append(p.items, item)
	return item, nil
}

func (p *precioRepoFull) ObtenerPrecio(id int) (*models.PrecioRecurso, error) {
	for _, item := range p.items {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, errors.New("not found")
}

func (p *precioRepoFull) ListarPrecios() []*models.PrecioRecurso {
	if p.items == nil {
		return []*models.PrecioRecurso{}
	}
	return p.items
}

func (p *precioRepoFull) ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	for i, item := range p.items {
		if item.ID == id {
			p.items[i].Precio = in.Precio
			p.items[i].FechaVigencia = in.FechaVigencia
			return p.items[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (p *precioRepoFull) EliminarPrecio(id int) error {
	for i, item := range p.items {
		if item.ID == id {
			p.items = append(p.items[:i], p.items[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (p *precioRepoFull) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	return p.items
}

func (p *precioRepoFull) PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	if len(p.items) > 0 {
		return p.items[len(p.items)-1], nil
	}
	return nil, errors.New("not found")
}

func (p *precioRepoFull) ExisteRecurso(tipo string, id int) error {
	return nil
}

func TestPrecioService_CRUD(t *testing.T) {
	svc := NewPreciosService(&precioRepoFull{})

	t.Run("Listado vacio", func(t *testing.T) {
		assert.Empty(t, svc.ListarPr())
	})

	t.Run("Obtener inexistente -> error", func(t *testing.T) {
		_, err := svc.ObtenerPr(1)
		assert.Error(t, err)
	})

	t.Run("Crear", func(t *testing.T) {
		in := models.EntradaPrecioRecurso{
			RecursoTipo:   "material",
			RecursoID:     1,
			Precio:        decimal.NewFromFloat(12.50),
			FechaVigencia: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		creado, err := svc.CrearPr(in)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromFloat(12.50).Equal(creado.Precio))
	})

	t.Run("Historial y vigente", func(t *testing.T) {
		hist := svc.HistorialPr("material", 1)
		assert.Len(t, hist, 1)

		vigente, err := svc.PrecioVigentePr("material", 1)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromFloat(12.50).Equal(vigente.Precio))
	})

	t.Run("Actualizar", func(t *testing.T) {
		in := models.EntradaPrecioRecurso{
			RecursoTipo:   "material",
			RecursoID:     1,
			Precio:        decimal.NewFromFloat(15.00),
			FechaVigencia: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		}
		actualizado, err := svc.ActualizarPr(1, in)
		require.NoError(t, err)
		assert.True(t, decimal.NewFromFloat(15.00).Equal(actualizado.Precio))
	})

	t.Run("Eliminar", func(t *testing.T) {
		err := svc.EliminarPr(1)
		assert.NoError(t, err)

		err = svc.EliminarPr(1)
		assert.Error(t, err)
	})
}
