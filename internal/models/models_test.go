package models

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEntradaMaterial_Validar(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		e := EntradaMaterial{Nombre: "Cemento", Unidad: "kg", PrecioReferencia: "10.50"}
		assert.NoError(t, e.ValidarMaterial())
	})
	t.Run("nombre vacio", func(t *testing.T) {
		e := EntradaMaterial{Nombre: "", Unidad: "kg", PrecioReferencia: "10.50"}
		assert.Error(t, e.ValidarMaterial())
	})
	t.Run("unidad no permitida", func(t *testing.T) {
		e := EntradaMaterial{Nombre: "Cemento", Unidad: "litro", PrecioReferencia: "10.50"}
		assert.Error(t, e.ValidarMaterial())
	})
	t.Run("precio cero", func(t *testing.T) {
		e := EntradaMaterial{Nombre: "Cemento", Unidad: "kg", PrecioReferencia: "0"}
		assert.Error(t, e.ValidarMaterial())
	})
}

func TestEntradaManoObra_Validar(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		e := EntradaManoObra{Descripcion: "Albañil", Categoria: "oficial", Unidad: "hora", CostoReferencia: decimal.RequireFromString("15.00")}
		assert.NoError(t, e.ValidarManoObra())
	})
	t.Run("categoria invalida", func(t *testing.T) {
		e := EntradaManoObra{Descripcion: "Albañil", Categoria: "gerente", Unidad: "hora", CostoReferencia: decimal.RequireFromString("15.00")}
		assert.Error(t, e.ValidarManoObra())
	})
}

func TestEntradaEquipo_Validar(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		e := EntradaEquipo{Nombre: "Excavadora", Tipo: "pesado", Unidad: "hora", CostoHora: decimal.RequireFromString("50.00")}
		assert.NoError(t, e.ValidarEquipo())
	})
	t.Run("tipo invalido", func(t *testing.T) {
		e := EntradaEquipo{Nombre: "Excavadora", Tipo: "mediano", Unidad: "hora", CostoHora: decimal.RequireFromString("50.00")}
		assert.Error(t, e.ValidarEquipo())
	})
}

func TestEntradaPrecioRecurso_Validar(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		e := EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 1, Precio: decimal.RequireFromString("10.00"), FechaVigencia: mustParseTime("2026-01-01")}
		assert.NoError(t, e.ValidarPrecio())
	})
	t.Run("tipo invalido", func(t *testing.T) {
		e := EntradaPrecioRecurso{RecursoTipo: "invalido", RecursoID: 1, Precio: decimal.RequireFromString("10.00"), FechaVigencia: mustParseTime("2026-01-01")}
		assert.Error(t, e.ValidarPrecio())
	})
	t.Run("precio cero", func(t *testing.T) {
		e := EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 1, Precio: decimal.Zero, FechaVigencia: mustParseTime("2026-01-01")}
		assert.Error(t, e.ValidarPrecio())
	})
}

func TestEntradaUsuario_Validar(t *testing.T) {
	t.Run("valido", func(t *testing.T) {
		e := EntradaUsuario{Email: "test@test.com", Password: "123456"}
		assert.NoError(t, e.ValidarUsuario())
	})
	t.Run("email vacio", func(t *testing.T) {
		e := EntradaUsuario{Email: "", Password: "123456"}
		assert.Error(t, e.ValidarUsuario())
	})
	t.Run("password corta", func(t *testing.T) {
		e := EntradaUsuario{Email: "test@test.com", Password: "12345"}
		assert.Error(t, e.ValidarUsuario())
	})
	t.Run("rol invalido", func(t *testing.T) {
		e := EntradaUsuario{Email: "test@test.com", Password: "123456", Rol: "superadmin"}
		assert.Error(t, e.ValidarUsuario())
	})
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
