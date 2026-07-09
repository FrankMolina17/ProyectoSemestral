package services

import (

	"testing"
	

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

type mockManoObraRepo struct {
	crearCalled bool
}

func (m *mockManoObraRepo) CrearManoObra(in models.EntradaManoObra) (*models.ManoObra, error) {
	m.crearCalled = true
	return &models.ManoObra{}, nil
}
func (m *mockManoObraRepo) ObtenerManoObra(id int) (*models.ManoObra, bool) { return nil, false }
func (m *mockManoObraRepo) ListarManoObra() []*models.ManoObra            { return nil }
func (m *mockManoObraRepo) ActualizarManoObra(id int, in models.EntradaManoObra) (*models.ManoObra, bool) {
	return nil, false
}
func (m *mockManoObraRepo) EliminarManoObra(id int) bool { return false }

func TestCrearManoObra_RechazaDatoInvalido_NoLlamaAlRepositorio(t *testing.T) {
	casos := []struct {
		nombre string
		in     models.EntradaManoObra
	}{
		{"descripcion vacia", models.EntradaManoObra{Descripcion: "", Categoria: "oficial", Unidad: "día", CostoReferencia: decimal.RequireFromString("10.00")}},
		{"categoria vacia", models.EntradaManoObra{Descripcion: "Oficial", Categoria: "", Unidad: "día", CostoReferencia: decimal.RequireFromString("10.00")}},
		{"unidad vacia", models.EntradaManoObra{Descripcion: "Oficial", Categoria: "oficial", Unidad: "", CostoReferencia: decimal.RequireFromString("10.00")}},
		{"costo cero", models.EntradaManoObra{Descripcion: "Oficial", Categoria: "oficial", Unidad: "día", CostoReferencia: decimal.Zero}},
		{"costo negativo", models.EntradaManoObra{Descripcion: "Oficial", Categoria: "oficial", Unidad: "día", CostoReferencia: decimal.RequireFromString("-5.00")}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			mock := &mockManoObraRepo{}
			svc := NewManoObraService(mock)

			_, err := svc.CrearMa(tc.in)

			assert.Error(t, err)
			assert.False(t, mock.crearCalled, "CrearManoObra no debe ser llamado cuando la entrada es invalida")
		})
	}
}

func TestCrearManoObra_DatoValido_LlamaAlRepositorio(t *testing.T) {
	mock := &mockManoObraRepo{}
	svc := NewManoObraService(mock)

	in := models.EntradaManoObra{
		Descripcion:     "Maestro de obra",
		Categoria:       "oficial",
		Unidad:          "día",
		CostoReferencia: decimal.RequireFromString("35.00"),
	}

	_, err := svc.CrearMa(in)

	assert.NoError(t, err)
	assert.True(t, mock.crearCalled, "CrearManoObra debe ser llamado cuando la entrada es valida")
}
