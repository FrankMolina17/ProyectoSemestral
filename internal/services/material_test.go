package services

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockMaterialRepo struct {
	crearCalled bool
}

func (m *mockMaterialRepo) CrearMateriales(in models.EntradaMaterial) (*models.Material, error) {
	m.crearCalled = true
	return nil, nil
}

func (m *mockMaterialRepo) ObtenerMateriales(id int) (*models.Material, bool) {
	return nil, false
}

func (m *mockMaterialRepo) ListarMateriales() []*models.Material {
	return nil
}

func (m *mockMaterialRepo) ActualizarMateriales(id int, in models.EntradaMaterial) (*models.Material, bool) {
	return nil, false
}

func (m *mockMaterialRepo) EliminarMateriales(id int) bool {
	return false
}

//TestCrearMaterial_RechazaDatoInvalido_NoLlamaAlRepositorio verifica que el servicio no llame al repositorio cuando los datos son inválidos.
func TestCrearMaterial_RechazaDatoInvalido_NoLlamaAlRepositorio(t *testing.T) {
	casos := []struct {
		nombre string
		in     models.EntradaMaterial
	}{
		{
			"nombre vacio",
			models.EntradaMaterial{
				Nombre:           "",
				Unidad:           "unidad",
				PrecioReferencia: decimal.NewFromFloat(10.0),
			},
		},
		{
			"unidad no permitida",
			models.EntradaMaterial{
				Nombre:           "Cemento",
				Unidad:           "km",
				PrecioReferencia: decimal.NewFromFloat(10.0),
			},
		},
		{
			"precio cero",
			models.EntradaMaterial{
				Nombre:           "Cemento",
				Unidad:           "unidad",
				PrecioReferencia: decimal.Zero,
			},
		},
		{
			"precio negativo",
			models.EntradaMaterial{
				Nombre:           "Cemento",
				Unidad:           "unidad",
				PrecioReferencia: decimal.NewFromFloat(-5.0),
			},
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			mock := &mockMaterialRepo{}
			svc := NewMaterialService(mock)

			_, err := svc.CrearM(tc.in)

			assert.Error(t, err)
			assert.False(t, mock.crearCalled, "CrearMateriales no debe ser llamado cuando la entrada es invalida")
		})
	}
}

// TestCrearMaterial_DatoValido_LlamaAlRepositorio verifica que el servicio llame al repositorio cuando los datos son válidos.
func TestCrearMaterial_DatoValido_LlamaAlRepositorio(t *testing.T) {
	mock := &mockMaterialRepo{}
	svc := NewMaterialService(mock)

	in := models.EntradaMaterial{
		Nombre:           "Cemento Portland",
		Descripcion:      "Saco de 50kg",
		Unidad:           "unidad",
		PrecioReferencia: decimal.NewFromFloat(25.50),
	}

	_, err := svc.CrearM(in)

	assert.NoError(t, err)
	assert.True(t, mock.crearCalled, "CrearMateriales debe ser llamado cuando la entrada es valida")
}
