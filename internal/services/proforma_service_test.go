package services_test

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProformaRepo registra si se llamó a CrearProforma (no persiste nada).
type mockProformaRepo struct {
	crearLlamado bool
}

func (m *mockProformaRepo) CrearProforma(p models.Proforma) (models.Proforma, error) {
	m.crearLlamado = true
	return p, nil
}

func (m *mockProformaRepo) ObtenerPorID(id int) (models.Proforma, error) {
	return models.Proforma{}, nil
}

func (m *mockProformaRepo) ObtenerTodos() ([]models.Proforma, error) {
	return nil, nil
}

func (m *mockProformaRepo) ActualizarProforma(id int, datos models.Proforma) (models.Proforma, error) {
	return models.Proforma{}, nil
}

func (m *mockProformaRepo) EliminarProforma(id int) error {
	return nil
}

func (m *mockProformaRepo) AprobarProforma(id int) (models.Proforma, error) {
	return models.Proforma{}, nil
}

func (m *mockProformaRepo) AgregarItem(item models.ProformaItem) (models.ProformaItem, error) {
	return models.ProformaItem{}, nil
}

func (m *mockProformaRepo) ObtenerItems(proformaID int) ([]models.ProformaItem, error) {
	return nil, nil
}

func (m *mockProformaRepo) AgregarNota(nota models.NotaProforma) (models.NotaProforma, error) {
	return models.NotaProforma{}, nil
}

func (m *mockProformaRepo) ObtenerNotas(proformaID int) ([]models.NotaProforma, error) {
	return nil, nil
}

func (m *mockProformaRepo) CrearCliente(c models.Cliente) (models.Cliente, error) {
	return models.Cliente{}, nil
}

func (m *mockProformaRepo) ObtenerClientes() ([]models.Cliente, error) {
	return nil, nil
}

func (m *mockProformaRepo) ObtenerClientePorID(id int) (models.Cliente, error) {
	return models.Cliente{}, nil
}

func (m *mockProformaRepo) ActualizarCliente(id int, datos models.Cliente) (models.Cliente, error) {
	return models.Cliente{}, nil
}

func (m *mockProformaRepo) EliminarCliente(id int) error {
	return nil
}

// TestCrearProforma_RechazaSinObraID verifica que una proforma sin obra_id
// sea rechazada por el service y NO llegue al repositorio.
func TestCrearProforma_RechazaSinObraID(t *testing.T) {
	mock := &mockProformaRepo{}
	svc := services.NuevoProformaService(mock)

	_, err := svc.CrearProforma(models.Proforma{
		Nombre: "Proforma edificio central",
		ObraID: 0,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrObraIDRequerido)
	assert.False(t, mock.crearLlamado, "el repositorio no debe ser invocado con datos inválidos")
}
