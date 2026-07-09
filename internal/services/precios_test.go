package services

import (

	"testing"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)


// ─────────────────────────────────────────────
// PRECIO
// ─────────────────────────────────────────────

type mockPrecioRepo struct {
	crearCalled bool
}

func (m *mockPrecioRepo) ListarPrecios() []*models.PrecioRecurso { return nil }
func (m *mockPrecioRepo) ObtenerPrecio(id int) (*models.PrecioRecurso, error) {
	return nil, nil
}
func (m *mockPrecioRepo) CrearPrecio(in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	m.crearCalled = true
	return &models.PrecioRecurso{}, nil
}
func (m *mockPrecioRepo) HistorialPrecios(tipo string, recursoID int) []*models.PrecioRecurso {
	return nil
}
func (m *mockPrecioRepo) PrecioVigente(tipo string, recursoID int) (*models.PrecioRecurso, error) {
	return nil, nil
}
func (m *mockPrecioRepo) ActualizarPrecio(id int, in models.EntradaPrecioRecurso) (*models.PrecioRecurso, error) {
	return nil, nil
}
func (m *mockPrecioRepo) ExisteRecurso(tipo string, id int) error        { return nil }
func (m *mockPrecioRepo) EliminarPrecio(id int) error                   { return nil }

func TestCrearPrecio_RechazaDatoInvalido_NoLlamaAlRepositorio(t *testing.T) {
	casos := []struct {
		nombre string
		in     models.EntradaPrecioRecurso
	}{
		{"recurso_tipo invalido", models.EntradaPrecioRecurso{RecursoTipo: "otro", RecursoID: 1, Precio: decimal.RequireFromString("10.00"), FechaVigencia: time.Now()}},
		{"recurso_id cero", models.EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 0, Precio: decimal.RequireFromString("10.00"), FechaVigencia: time.Now()}},
		{"precio cero", models.EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 1, Precio: decimal.Zero, FechaVigencia: time.Now()}},
		{"precio negativo", models.EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 1, Precio: decimal.RequireFromString("-5.00"), FechaVigencia: time.Now()}},
		{"fecha vigencia vacia", models.EntradaPrecioRecurso{RecursoTipo: "material", RecursoID: 1, Precio: decimal.RequireFromString("10.00"), FechaVigencia: time.Time{}}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			mock := &mockPrecioRepo{}
			svc := NewPreciosService(mock)

			_, err := svc.CrearPr(tc.in)

			assert.Error(t, err)
			assert.False(t, mock.crearCalled, "CrearPrecio no debe ser llamado cuando la entrada es invalida")
		})
	}
}

func TestCrearPrecio_DatoValido_LlamaAlRepositorio(t *testing.T) {
	mock := &mockPrecioRepo{}
	svc := NewPreciosService(mock)

	in := models.EntradaPrecioRecurso{
		RecursoTipo:   "material",
		RecursoID:     1,
		Precio:        decimal.RequireFromString("25.50"),
		FechaVigencia: time.Now(),
	}

	_, err := svc.CrearPr(in)

	assert.NoError(t, err)
	assert.True(t, mock.crearCalled, "CrearPrecio debe ser llamado cuando la entrada es valida")
}