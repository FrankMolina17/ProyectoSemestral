package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
)


func RepuestaJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

//su funcion es devolver un 200
func ok(w http.ResponseWriter, data any) {
	RepuestaJSON(w, http.StatusOK, map[string]any{"data": data})
}

func creando(w http.ResponseWriter, data any, id int) {
	RepuestaJSON(w, http.StatusCreated, map[string]any{"data": data, "id": id})
}

//su funcion es devolver un 400
func MalFormado(w http.ResponseWriter, msg string) {
	RepuestaJSON(w, http.StatusBadRequest, map[string]string{"error": msg})
}
 
func NoEncontrado(w http.ResponseWriter, recurso string, id int) {
	RepuestaJSON(w, http.StatusNotFound,
		map[string]string{"error": fmt.Sprintf("%s con id %d no encontrado", recurso, id)})
}

func ErrorMermoria(w http.ResponseWriter, err error, recurso string, id int) {
	switch {
	case errors.Is(err, storage.ErrNotFound):
		NoEncontrado(w, recurso, id)
	case errors.Is(err, storage.ErrDuplicated):
		RepuestaJSON(w, http.StatusConflict,
			map[string]string{"error": "nombre ya existe para esa unidad"})
	default:
		RepuestaJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "error interno"})
	}
}

//esto es para decodificar el json y mostrar el error
func DecodificarJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		MalFormado(w, "body malformado: "+err.Error())
		return false
	}
	return true
}

//esto es para obtener el id
func ParaObtenerelID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		MalFormado(w, "id debe ser un entero positivo")
		return 0, false
	}
	return id, true
}

// ─────────────────────────────────────────────
// MATERIAL
// ─────────────────────────────────────────────

type MaterialHandler struct{
	 s *storage.Storage 
}

func NewMaterialHandler(s *storage.Storage) *MaterialHandler { 
	return &MaterialHandler{s} 
}

// GET /material  →  200
func (h *MaterialHandler) ListandoMateriales(w http.ResponseWriter, r *http.Request) {
	ok(w, h.s.ListarMateriales())
}

// GET /material/{id}  →  200 | 404
func (h *MaterialHandler) ObtenerMaterialPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	mat, err := h.s.ObtenerMateriales(id)
	if err != nil {
		ErrorMermoria(w, err, "material", id)
		return
	}
	ok(w, mat)
}

// POST /material  →  201 | 400
func (h *MaterialHandler) CreandoMaterial(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaMaterial
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	mat, err := h.s.CrearMateriales(in)
	if err != nil {
		ErrorMermoria(w, err, "material", 0)
		return
	}
	creando(w, mat, mat.ID)
}

// PUT /material/{id}  →  200 | 400 | 404
func (h *MaterialHandler) ActulizarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaMaterial
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	mat, err := h.s.ActualizarMateriales(id, in)
	if err != nil {
		ErrorMermoria(w, err, "material", id)
		return
	}
	ok(w, mat)
}

// DELETE /material/{id}  →  204 | 404
func (h *MaterialHandler) BorrarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.s.EliminarMateriales(id); err != nil {
		ErrorMermoria(w, err, "material", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

type ManoObraHandler struct{ 
	s *storage.Storage 
}

func NewManoObraHandler(s *storage.Storage) *ManoObraHandler { 
	return &ManoObraHandler{s} 
}

// GET /manoobra  →  200
func (h *ManoObraHandler) ListarUnaManoObra(w http.ResponseWriter, r *http.Request) {
	ok(w, h.s.ListarManoObra())
}

// GET /manoobra/{id}  →  200 | 404
func (h *ManoObraHandler) ObtenerUnaManoObraPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	mo, err := h.s.ObtenerManoObra(id)
	if err != nil {
		ErrorMermoria(w, err, "mano_obra", id)
		return
	}
	ok(w, mo)
}

// POST /manoobra  →  201 | 400
func (h *ManoObraHandler) CreandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaManoObra
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	mo, err := h.s.CrearManoObra(in)
	if err != nil {
		ErrorMermoria(w, err, "mano_obra", 0)
		return
	}
	creando(w, mo, mo.ID)
}

// PUT /manoobra/{id}  →  200 | 400 | 404
func (h *ManoObraHandler) ActualizadoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaManoObra
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	mo, err := h.s.ActualizarManoObra(id, in)
	if err != nil {
		ErrorMermoria(w, err, "mano_obra", id)
		return
	}
	ok(w, mo)
}

// DELETE /manoobra/{id}  →  204 | 404
func (h *ManoObraHandler) BorrandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.s.EliminarManoObra(id); err != nil {
		ErrorMermoria(w, err, "mano_obra", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

type EquipoHandler struct{ 
	s *storage.Storage 
}

func NewEquipoHandler(s *storage.Storage) *EquipoHandler { 
	return &EquipoHandler{s} 
}

// GET /equipo  →  200
func (h *EquipoHandler) ListandoLosEquipos(w http.ResponseWriter, r *http.Request) {
	ok(w, h.s.ListarEquipos())
}

// GET /equipo/{id}  →  200 | 404
func (h *EquipoHandler) ObtenerUnEquipoPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	eq, err := h.s.ObtenerEquipo(id)
	if err != nil {
		ErrorMermoria(w, err, "equipo", id)
		return
	}
	ok(w, eq)
}

// POST /equipo  →  201 | 400
func (h *EquipoHandler) CrearUnEquipo(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaEquipo
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	eq, err := h.s.CrearEquipo(in)
	if err != nil {
		ErrorMermoria(w, err, "equipo", 0)
		return
	}
	creando(w, eq, eq.ID)
}

// PUT /equipo/{id}  →  200 | 400 | 404
func (h *EquipoHandler) ActualizarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaEquipo
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	eq, err := h.s.ActualizarEquipo(id, in)
	if err != nil {
		ErrorMermoria(w, err, "equipo", id)
		return
	}
	ok(w, eq)
}

// DELETE /equipo/{id}  →  204 | 404
func (h *EquipoHandler) BorrarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.s.EliminarEquipo(id); err != nil {
		ErrorMermoria(w, err, "equipo", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// PRECIO
// ─────────────────────────────────────────────

type PrecioHandler struct{ 
	s *storage.Storage 
}

func NewPrecioHandler(s *storage.Storage) *PrecioHandler { return &PrecioHandler{s} }

// GET /precio  →  200
func (h *PrecioHandler) ListarDeLosPrecios(w http.ResponseWriter, r *http.Request) {
	ok(w, h.s.ListarPrecios())
}

// POST /precio  →  201 | 400
func (h *PrecioHandler) CrearUnPrecio(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaPrecioRecurso
	if !DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		MalFormado(w, err.Error())
		return
	}
	p, err := h.s.CrearPrecio(in)
	if err != nil {
		ErrorMermoria(w, err, "precio", 0)
		return
	}
	creando(w, p, p.ID)
}

// GET /precio/{id}  →  200 | 404 

func (h *PrecioHandler) ObtenerUnPrecioPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := ParaObtenerelID(w, r)
	if !valid {
		return
	}
	p, err := h.s.ObtenerPrecio(id)
	if err != nil {
		ErrorMermoria(w, err, "precio", id)
		return
	}
	ok(w, p)
}


