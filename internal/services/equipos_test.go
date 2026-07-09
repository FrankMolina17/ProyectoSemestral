package services

import (
	
	"testing"
	

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

type mockEquipoRepo struct {
	crearCalled bool
}

func (m *mockEquipoRepo) CrearEquipo(in models.EntradaEquipo) (*models.Equipo, error) {
	m.crearCalled = true
	return &models.Equipo{}, nil
}
func (m *mockEquipoRepo) ObtenerEquipo(id int) (*models.Equipo, error) { return nil, nil }
func (m *mockEquipoRepo) ListarEquipos() []*models.Equipo              { return nil }
func (m *mockEquipoRepo) ActualizarEquipo(id int, in models.EntradaEquipo) (*models.Equipo, error) {
	return nil, nil
}
func (m *mockEquipoRepo) EliminarEquipo(id int) error { return nil }

func TestCrearEquipo_RechazaDatoInvalido_NoLlamaAlRepositorio(t *testing.T) {
	casos := []struct {
		nombre string
		in     models.EntradaEquipo
	}{
		{"nombre vacio", models.EntradaEquipo{Nombre: "", Tipo: "pesado", Unidad: "hora", CostoHora: decimal.RequireFromString("10.00"), Disponible: true}},
		{"tipo invalido", models.EntradaEquipo{Nombre: "Excavadora", Tipo: "volador", Unidad: "hora", CostoHora: decimal.RequireFromString("10.00"), Disponible: true}},
		{"unidad vacia", models.EntradaEquipo{Nombre: "Excavadora", Tipo: "pesado", Unidad: "", CostoHora: decimal.RequireFromString("10.00"), Disponible: true}},
		{"costo cero", models.EntradaEquipo{Nombre: "Excavadora", Tipo: "pesado", Unidad: "hora", CostoHora: decimal.Zero, Disponible: true}},
		{"costo negativo", models.EntradaEquipo{Nombre: "Excavadora", Tipo: "pesado", Unidad: "hora", CostoHora: decimal.RequireFromString("-5.00"), Disponible: true}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			mock := &mockEquipoRepo{}
			svc := NewEquipoService(mock)

			_, err := svc.CrearE(tc.in)

			assert.Error(t, err)
			assert.False(t, mock.crearCalled, "CrearEquipo no debe ser llamado cuando la entrada es invalida")
		})
	}
}

func TestCrearEquipo_DatoValido_LlamaAlRepositorio(t *testing.T) {
	mock := &mockEquipoRepo{}
	svc := NewEquipoService(mock)

	in := models.EntradaEquipo{
		Nombre:     "Excavadora CAT 320",
		Tipo:       "pesado",
		Unidad:     "hora",
		CostoHora:  decimal.RequireFromString("85.00"),
		Disponible: true,
	}

	_, err := svc.CrearE(in)

	assert.NoError(t, err)
	assert.True(t, mock.crearCalled, "CrearEquipo debe ser llamado cuando la entrada es valida")
}
