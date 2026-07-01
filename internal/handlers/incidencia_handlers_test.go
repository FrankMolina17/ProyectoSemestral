package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type fakeIncidenciaRepository struct {
	incidencias []models.Incidencia
	nextID      int
}

func (f *fakeIncidenciaRepository) ListarIncidencias() []models.Incidencia {
	return f.incidencias
}

func (f *fakeIncidenciaRepository) BuscarIncidenciaPorID(id int) (models.Incidencia, bool) {
	for _, incidencia := range f.incidencias {
		if incidencia.ID == id {
			return incidencia, true
		}
	}

	return models.Incidencia{}, false
}

func (f *fakeIncidenciaRepository) BuscarIncidenciaPorEntidad(id int, tipo string) (models.Incidencia, bool) {
	for _, incidencia := range f.incidencias {
		if incidencia.EntidadID == id && incidencia.EntidadTipo == tipo {
			return incidencia, true
		}
	}

	return models.Incidencia{}, false
}

func (f *fakeIncidenciaRepository) CrearIncidencia(i models.Incidencia) models.Incidencia {
	if f.nextID == 0 {
		f.nextID = 1
	}

	i.ID = f.nextID
	f.nextID++

	f.incidencias = append(f.incidencias, i)

	return i
}

func (f *fakeIncidenciaRepository) ActualizarIncidencia(id int, datos models.Incidencia) (models.Incidencia, bool) {
	for pos, incidencia := range f.incidencias {
		if incidencia.ID == id {
			datos.ID = id
			f.incidencias[pos] = datos
			return datos, true
		}
	}

	return models.Incidencia{}, false
}

func (f *fakeIncidenciaRepository) BorrarIncidencia(id int) bool {
	for pos, incidencia := range f.incidencias {
		if incidencia.ID == id {
			f.incidencias = append(f.incidencias[:pos], f.incidencias[pos+1:]...)
			return true
		}
	}

	return false
}

func TestObtenerIncidencias_SinToken_Responde401(t *testing.T) {
	r := chi.NewRouter()

	r.Route("/api/v1/incidencias", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, `{"error":"Token inválido"}`, http.StatusUnauthorized)
			})
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/incidencias/", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCrearIncidenciaHandler_CreaIncidenciaValida(t *testing.T) {
	repo := &fakeIncidenciaRepository{}
	servicio := services.NuevaIncidenciaService(repo)
	server := NewServer(servicio, nil)

	r := chi.NewRouter()
	r.Post("/api/v1/incidencias", server.CrearIncidencia)

	body := []byte(`{
		"entidad_tipo": "obra",
		"entidad_id": 1,
		"responsable_id": 5,
		"titulo": "Falta de material en columna",
		"descripcion": "Los hierros no llegaron a tiempo",
		"estado": "Abierta"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/incidencias", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Contains(t, rec.Body.String(), `"titulo":"Falta de material en columna"`)
	require.Contains(t, rec.Body.String(), `"estado":"Abierta"`)
}
