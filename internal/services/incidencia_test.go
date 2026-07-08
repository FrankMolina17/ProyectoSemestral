package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"errors"
	"testing"
)

type mockIncidenciaRepository struct {
	shouldReturnErrorOnUpdate bool
	updateWasCalled           bool
	updateIDPassed            int
	updateIncidencePassed     models.Incidencia
}

func (m *mockIncidenciaRepository) CrearIncidencia(c models.Incidencia) models.Incidencia {

	return c
}

func (m *mockIncidenciaRepository) ListarIncidencias() []models.Incidencia {

	return nil
}

func (m *mockIncidenciaRepository) BuscarIncidenciaPorID(id int) (models.Incidencia, bool) {

	return models.Incidencia{}, false
}

func (m *mockIncidenciaRepository) BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool) {

	return models.Incidencia{}, false
}

func (m *mockIncidenciaRepository) ActualizarIncidencia(id int, c models.Incidencia) (models.Incidencia, bool) {

	m.updateWasCalled = true
	m.updateIDPassed = id
	m.updateIncidencePassed = c

	if m.shouldReturnErrorOnUpdate {
		return models.Incidencia{}, false
	}

	return c, true
}

func (m *mockIncidenciaRepository) BorrarIncidencia(id int) bool {
	return true
}

func TestIncidenciaService_ActualizarIncidencia_InvalidTitle(t *testing.T) {

	mockRepo := &mockIncidenciaRepository{}
	service := NuevaIncidenciaService(mockRepo)

	invalidIncidencia := models.Incidencia{
		ID:          99,
		Titulo:      "",
		Descripcion: "Aa",
		Estado:      "Abierta",
	}

	idToUpdate := 99

	_, err := service.ActualizarIncidencia(idToUpdate, invalidIncidencia)

	if err == nil {
		t.Fatalf("Se esperaba un error por un título vacío, pero se obtuvo nil")
	}
	if !errors.Is(err, ErrTituloIncidenciaVacio) && err.Error() != ErrTituloIncidenciaVacio.Error() {
		t.Errorf("Error esperado '%v', Se obtuvo '%v'", ErrTituloIncidenciaVacio, err)
	}
	if mockRepo.updateWasCalled {
		t.Errorf("Se esperaba que el método «ActualizarIncidencia» del repositorio NO se llamara, pero se llamó con ID: %d, Incidencia: %+v", mockRepo.updateIDPassed, mockRepo.updateIncidencePassed)
	}
}
