package handlers

import (
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

// ─────────────────────────────────────────────
// MATERIAL
// ─────────────────────────────────────────────

type MaterialHandler struct {
	s *storage.Storage
}

func NewMaterialHandler(s *storage.Storage) *MaterialHandler {
	return &MaterialHandler{s}
}

// GET /material  →  200
func (h *MaterialHandler) ListandoMateriales(w http.ResponseWriter, r *http.Request) {
	material := h.s.ListarMateriales()
	RespondJSON(w, http.StatusOK, material)
}

// GET /material/{id}  →  200 | 404
func (h *MaterialHandler) ObtenerMaterialPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	material, ok := h.s.ObtenerMateriales(id)
	if !ok {
		services.NoEncontrado(w, "material", id)
		return
	}
	services.Ok(w, material)
}

// POST /material  →  201 | 400
func (h *MaterialHandler) CreandoMaterial(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaMaterial
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarMaterial(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	mat, err := h.s.CrearMateriales(in)
	if err != nil {
		services.ErrorMermoria(w, err, "material", 0)
		return
	}
	services.Creando(w, mat, mat.ID)
}

// PUT /material/{id}  →  200 | 400 | 404
func (h *MaterialHandler) ActulizarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaMaterial
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarMaterial(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	mat, ok := h.s.ActualizarMateriales(id, in)
	if !ok {
		services.NoEncontrado(w, "material", id)
		return
	}
	services.Ok(w, mat)
}

// DELETE /material/{id}  →  204 | 404
func (h *MaterialHandler) BorrarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if !h.s.EliminarMateriales(id) {
		services.NoEncontrado(w, "material", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

type ManoObraHandler struct {
	s *storage.Storage
}

func NewManoObraHandler(s *storage.Storage) *ManoObraHandler {
	return &ManoObraHandler{s}
}

// GET /manoobra  →  200
func (h *ManoObraHandler) ListarUnaManoObra(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.s.ListarManoObra())
}

// GET /manoobra/{id}  →  200 | 404
func (h *ManoObraHandler) ObtenerUnaManoObraPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	mo, ok := h.s.ObtenerManoObra(id)
	if !ok {
		services.NoEncontrado(w, "mano_obra", id)
		return
	}
	services.Ok(w, mo)
}

// POST /manoobra  →  201 | 400
func (h *ManoObraHandler) CreandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaManoObra
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarManoObra(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	mo, err := h.s.CrearManoObra(in)
	if err != nil {
		services.ErrorMermoria(w, err, "mano_obra", 0)
		return
	}
	services.Creando(w, mo, mo.ID)
}

// PUT /manoobra/{id}  →  200 | 400 | 404
func (h *ManoObraHandler) ActualizadoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaManoObra
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarManoObra(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	mo, ok := h.s.ActualizarManoObra(id, in)
	if !ok {
		services.NoEncontrado(w, "mano_obra", id)
		return
	}
	services.Ok(w, mo)
}

// DELETE /manoobra/{id}  →  204 | 404
func (h *ManoObraHandler) BorrandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if !h.s.EliminarManoObra(id) {
		services.NoEncontrado(w, "mano_obra", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

type EquipoHandler struct {
	s *storage.Storage
}

func NewEquipoHandler(s *storage.Storage) *EquipoHandler {
	return &EquipoHandler{s}
}

// GET /equipo  →  200
func (h *EquipoHandler) ListandoLosEquipos(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.s.ListarEquipos())
}

// GET /equipo/{id}  →  200 | 404
func (h *EquipoHandler) ObtenerUnEquipoPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	eq, err := h.s.ObtenerEquipo(id)
	if err != nil {
		services.ErrorMermoria(w, err, "equipo", id)
		return
	}
	services.Ok(w, eq)
}

// POST /equipo  →  201 | 400
func (h *EquipoHandler) CrearUnEquipo(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaEquipo
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarEquipo(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	eq, err := h.s.CrearEquipo(in)
	if err != nil {
		services.ErrorMermoria(w, err, "equipo", 0)
		return
	}
	services.Creando(w, eq, eq.ID)
}

// PUT /equipo/{id}  →  200 | 400 | 404
func (h *EquipoHandler) ActualizarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaEquipo
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarEquipo(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	eq, err := h.s.ActualizarEquipo(id, in)
	if err != nil {
		services.ErrorMermoria(w, err, "equipo", id)
		return
	}
	services.Ok(w, eq)
}

// DELETE /equipo/{id}  →  204 | 404
func (h *EquipoHandler) BorrarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.s.EliminarEquipo(id); err != nil {
		services.ErrorMermoria(w, err, "equipo", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// PRECIO
// ─────────────────────────────────────────────

type PrecioHandler struct {
	s *storage.Storage
}

func NewPrecioHandler(s *storage.Storage) *PrecioHandler { return &PrecioHandler{s} }

// GET /precio  →  200
func (h *PrecioHandler) ListarDeLosPrecios(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.s.ListarPrecios())
}

// POST /precio  →  201 | 400
func (h *PrecioHandler) CrearUnPrecio(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaPrecioRecurso
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarPrecio(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	p, err := h.s.CrearPrecio(in)
	if err != nil {
		services.ErrorMermoria(w, err, in.RecursoTipo, in.RecursoID)
		return
	}
	services.Creando(w, p, p.ID)
}

// GET /precio/{tipo}/{recursoID}  →  200
func (h *PrecioHandler) HistorialPorRecurso(w http.ResponseWriter, r *http.Request) {
	tipo, recursoID, valid := services.ParaObtenerTipoRecursoID(w, r)
	if !valid {
		return
	}
	services.Ok(w, h.s.HistorialPrecios(tipo, recursoID))
}

// GET /precio/{tipo}/{recursoID}/vigente  →  200 | 404
func (h *PrecioHandler) PrecioVigentePorRecurso(w http.ResponseWriter, r *http.Request) {
	tipo, recursoID, valid := services.ParaObtenerTipoRecursoID(w, r)
	if !valid {
		return
	}
	p, err := h.s.PrecioVigente(tipo, recursoID)
	if err != nil {
		services.ErrorMermoria(w, err, "precio vigente", recursoID)
		return
	}
	services.Ok(w, p)
}

// GET /precio/{id}  →  200 | 404
func (h *PrecioHandler) ObtenerUnPrecioPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	p, err := h.s.ObtenerPrecio(id)
	if err != nil {
		services.ErrorMermoria(w, err, "precio", id)
		return
	}
	services.Ok(w, p)
}

// PUT /precio/{id}  →  200 | 400 | 404
func (h *PrecioHandler) ActualizarUnPrecio(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	var in models.EntradaPrecioRecurso
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarPrecio(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	p, err := h.s.ActualizarPrecio(id, in)
	if err != nil {
		services.ErrorMermoria(w, err, "precio", id)
		return
	}
	services.Ok(w, p)
}

// DELETE /precio/{id}  →  204 | 404
func (h *PrecioHandler) BorrarUnPrecio(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.s.EliminarPrecio(id); err != nil {
		services.ErrorMermoria(w, err, "precio", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}