package storage

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMemoria(t *testing.T) *Storage {
	return New()
}

func TestMemoria_Materiales(t *testing.T) {
	s := setupMemoria(t)

	in := models.EntradaMaterial{Nombre: "Cemento", Unidad: "kg", PrecioReferencia: "10.50"}
	m, err := s.CrearMateriales(in)
	require.NoError(t, err)
	assert.Equal(t, "Cemento", m.Nombre)

	lista := s.ListarMateriales()
	assert.Len(t, lista, 1)

	obt, ok := s.ObtenerMateriales(m.ID)
	assert.True(t, ok)
	assert.Equal(t, "Cemento", obt.Nombre)

	act, ok := s.ActualizarMateriales(m.ID, models.EntradaMaterial{Nombre: "Cemento Plus", Unidad: "kg", PrecioReferencia: "12.00"})
	assert.True(t, ok)
	assert.Equal(t, "Cemento Plus", act.Nombre)

	assert.True(t, s.EliminarMateriales(m.ID))
	assert.False(t, s.EliminarMateriales(m.ID))

	_, ok = s.ObtenerMateriales(m.ID)
	assert.False(t, ok)

	_, ok = s.ObtenerMateriales(999)
	assert.False(t, ok)

	_, ok = s.ActualizarMateriales(999, in)
	assert.False(t, ok)
}

func TestMemoria_ManoObra(t *testing.T) {
	s := setupMemoria(t)

	in := models.EntradaManoObra{
		Descripcion:     "Albañil",
		Categoria:       "oficial",
		Unidad:          "hora",
		CostoReferencia: decimal.RequireFromString("15.00"),
	}
	m, err := s.CrearManoObra(in)
	require.NoError(t, err)
	assert.Equal(t, "Albañil", m.Descripcion)

	lista := s.ListarManoObra()
	assert.Len(t, lista, 1)

	obt, ok := s.ObtenerManoObra(m.ID)
	assert.True(t, ok)
	assert.Equal(t, "oficial", obt.Categoria)

	act, ok := s.ActualizarManoObra(m.ID, models.EntradaManoObra{
		Descripcion:     "Albañil Senior",
		Categoria:       "oficial",
		Unidad:          "hora",
		CostoReferencia: decimal.RequireFromString("20.00"),
	})
	assert.True(t, ok)
	assert.Equal(t, "Albañil Senior", act.Descripcion)

	assert.True(t, s.EliminarManoObra(m.ID))
	assert.False(t, s.EliminarManoObra(m.ID))

	_, ok = s.ObtenerManoObra(999)
	assert.False(t, ok)

	_, ok = s.ActualizarManoObra(999, in)
	assert.False(t, ok)
}

func TestMemoria_Equipos(t *testing.T) {
	s := setupMemoria(t)

	in := models.EntradaEquipo{
		Nombre:     "Excavadora",
		Tipo:       "pesado",
		Unidad:     "hora",
		CostoHora:  decimal.RequireFromString("85.00"),
		Disponible: true,
	}
	m, err := s.CrearEquipo(in)
	require.NoError(t, err)
	assert.Equal(t, "Excavadora", m.Nombre)

	lista := s.ListarEquipos()
	assert.Len(t, lista, 1)

	obt, err := s.ObtenerEquipo(m.ID)
	require.NoError(t, err)
	assert.Equal(t, "pesado", obt.Tipo)

	act, err := s.ActualizarEquipo(m.ID, models.EntradaEquipo{
		Nombre:     "Excavadora CAT",
		Tipo:       "pesado",
		Unidad:     "hora",
		CostoHora:  decimal.RequireFromString("90.00"),
		Disponible: false,
	})
	require.NoError(t, err)
	assert.Equal(t, "Excavadora CAT", act.Nombre)

	assert.NoError(t, s.EliminarEquipo(m.ID))
	assert.Error(t, s.EliminarEquipo(m.ID))

	_, err = s.ObtenerEquipo(m.ID)
	assert.Error(t, err)

	_, err = s.ObtenerEquipo(999)
	assert.Error(t, err)
}

func TestMemoria_Precios(t *testing.T) {
	s := setupMemoria(t)

	mat, err := s.CrearMateriales(models.EntradaMaterial{Nombre: "Arena", Unidad: "m³", PrecioReferencia: "22.00"})
	require.NoError(t, err)

	in := models.EntradaPrecioRecurso{
		RecursoTipo: "material",
		RecursoID:   mat.ID,
		Precio:      decimal.RequireFromString("22.00"),
	}
	p, err := s.CrearPrecio(in)
	require.NoError(t, err)
	assert.Greater(t, p.ID, 0)

	lista := s.ListarPrecios()
	assert.Len(t, lista, 1)

	obt, err := s.ObtenerPrecio(p.ID)
	require.NoError(t, err)
	assert.Equal(t, p.ID, obt.ID)

	assert.NoError(t, s.ExisteRecurso("material", mat.ID))
	assert.Error(t, s.ExisteRecurso("material", 99999))

	hist := s.HistorialPrecios("material", mat.ID)
	assert.Len(t, hist, 1)

	vig, err := s.PrecioVigente("material", mat.ID)
	require.NoError(t, err)
	assert.Equal(t, "22", vig.Precio.String())

	_, err = s.PrecioVigente("material", 99999)
	assert.Error(t, err)

	act, err := s.ActualizarPrecio(p.ID, models.EntradaPrecioRecurso{
		RecursoTipo: "material",
		RecursoID:   mat.ID,
		Precio:      decimal.RequireFromString("25.00"),
	})
	require.NoError(t, err)
	assert.Equal(t, "25", act.Precio.String())

	assert.NoError(t, s.EliminarPrecio(p.ID))
	_, err = s.ObtenerPrecio(p.ID)
	assert.Error(t, err)
}

func TestMemoria_Usuarios(t *testing.T) {
	s := setupMemoria(t)

	u, err := s.CrearUsuario(models.EntradaUsuario{
		Email:    "test@test.com",
		Password: "123456",
	})
	require.NoError(t, err)
	assert.Equal(t, "test@test.com", u.Email)

	lista := s.ListarUsuarios()
	assert.Len(t, lista, 1)

	obt, ok := s.ObtenerUsuarioPorID(u.ID)
	assert.True(t, ok)
	assert.Equal(t, "test@test.com", obt.Email)

	_, ok = s.ObtenerUsuarioPorID(999)
	assert.False(t, ok)

	usr, ok := s.BuscarUsuarioPorEmail("test@test.com")
	assert.True(t, ok)
	assert.Equal(t, u.ID, usr.ID)

	_, ok = s.BuscarUsuarioPorEmail("noexiste@test.com")
	assert.False(t, ok)
}

func TestMemoria_Seed(t *testing.T) {
	s := setupMemoria(t)
	s.Seed()

	mats := s.ListarMateriales()
	assert.NotEmpty(t, mats)

	manos := s.ListarManoObra()
	assert.NotEmpty(t, manos)

	equipos := s.ListarEquipos()
	assert.NotEmpty(t, equipos)

	precios := s.ListarPrecios()
	assert.NotEmpty(t, precios)
}
