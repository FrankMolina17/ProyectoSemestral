package handlers

import (
	"errors"
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

// ─────────────────────────────────────────────
// MATERIAL
// ─────────────────────────────────────────────

type MaterialHandler struct {
	svc *services.MaterialService
}

func NewMaterialHandler(svc *services.MaterialService) *MaterialHandler {
	return &MaterialHandler{svc: svc}
}

// GET /material  →  200
func (h *MaterialHandler) ListandoMateriales(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.svc.Listado())
}

// GET /material/{id}  →  200 | 404
func (h *MaterialHandler) ObtenerMaterialPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	m, err := h.svc.ObtenerM(id)
	if err != nil {
		notFoundOrError(w, err, "material", id)
		return
	}
	services.Ok(w, m)
}

// POST /material  →  201 | 400 | 409
func (h *MaterialHandler) CreandoMaterial(w http.ResponseWriter, r *http.Request) {
	var in models.EntradaMaterial
	if !services.DecodificarJSON(w, r, &in) {
		return
	}
	if err := in.ValidarMaterial(); err != nil {
		services.MalFormado(w, err.Error())
		return
	}
	m, err := h.svc.CrearM(in)
	if err != nil {
		services.ErrorMermoria(w, err, "material", 0)
		return
	}
	services.Creando(w, m, m.ID)
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
	m, err := h.svc.ActualizarM(id, in)
	if err != nil {
		notFoundOrError(w, err, "material", id)
		return
	}
	services.Ok(w, m)
}

// DELETE /material/{id}  →  204 | 404
func (h *MaterialHandler) BorrarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.svc.EliminarM(id); err != nil {
		notFoundOrError(w, err, "material", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// MANO DE OBRA
// ─────────────────────────────────────────────

type ManoObraHandler struct {
	svc *services.ManoObraServise
}

func NewManoObraHandler(svc *services.ManoObraServise) *ManoObraHandler {
	return &ManoObraHandler{svc: svc}
}

// GET /manoobra  →  200
func (h *ManoObraHandler) ListarUnaManoObra(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.svc.ListadoMa())
}

// GET /manoobra/{id}  →  200 | 404
func (h *ManoObraHandler) ObtenerUnaManoObraPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	m, err := h.svc.ObtenerMa(id)
	if err != nil {
		notFoundOrError(w, err, "mano_obra", id)
		return
	}
	services.Ok(w, m)
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
	m, err := h.svc.CrearMa(in)
	if err != nil {
		services.ErrorMermoria(w, err, "mano_obra", 0)
		return
	}
	services.Creando(w, m, m.ID)
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
	m, err := h.svc.ActualizarMa(id, in)
	if err != nil {
		notFoundOrError(w, err, "mano_obra", id)
		return
	}
	services.Ok(w, m)
}

// DELETE /manoobra/{id}  →  204 | 404
func (h *ManoObraHandler) BorrandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if !h.svc.EliminarMa(id) {
		services.NoEncontrado(w, "mano_obra", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// EQUIPO
// ─────────────────────────────────────────────

type EquipoHandler struct {
	svc *services.EquipoService
}

func NewEquipoHandler(svc *services.EquipoService) *EquipoHandler {
	return &EquipoHandler{svc: svc}
}

// GET /equipo  →  200
func (h *EquipoHandler) ListandoLosEquipos(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.svc.ListadoE())
}

// GET /equipo/{id}  →  200 | 404
func (h *EquipoHandler) ObtenerUnEquipoPorID(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	e, err := h.svc.ObtenerE(id)
	if err != nil {
		notFoundOrError(w, err, "equipo", id)
		return
	}
	services.Ok(w, e)
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
	e, err := h.svc.CrearE(in)
	if err != nil {
		services.ErrorMermoria(w, err, "equipo", 0)
		return
	}
	services.Creando(w, e, e.ID)
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
	e, err := h.svc.ActualizarE(id, in)
	if err != nil {
		notFoundOrError(w, err, "equipo", id)
		return
	}
	services.Ok(w, e)
}

// DELETE /equipo/{id}  →  204 | 404
func (h *EquipoHandler) BorrarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	if err := h.svc.EliminarE(id); err != nil {
		notFoundOrError(w, err, "equipo", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// PRECIO
// ─────────────────────────────────────────────

type PrecioHandler struct {
	svc *services.PreciosService
}

func NewPrecioHandler(svc *services.PreciosService) *PrecioHandler {
	return &PrecioHandler{svc: svc}
}

// GET /precio  →  200
func (h *PrecioHandler) ListarDeLosPrecios(w http.ResponseWriter, r *http.Request) {
	services.Ok(w, h.svc.ListarPr())
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
	p, err := h.svc.CrearPr(in)
	if err != nil {
		services.ErrorMermoria(w, err, "precio", 0)
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
	services.Ok(w, h.svc.HistorialPr(tipo, recursoID))
}

// GET /precio/{tipo}/{recursoID}/vigente  →  200 | 404
func (h *PrecioHandler) PrecioVigentePorRecurso(w http.ResponseWriter, r *http.Request) {
	tipo, recursoID, valid := services.ParaObtenerTipoRecursoID(w, r)
	if !valid {
		return
	}
	p, err := h.svc.PrecioVigentePr(tipo, recursoID)
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
	p, err := h.svc.ObtenerPr(id)
	if err != nil {
		notFoundOrError(w, err, "precio", id)
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
	p, err := h.svc.ActualizarPr(id, in)
	if err != nil {
		notFoundOrError(w, err, "precio", id)
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
	if err := h.svc.EliminarPr(id); err != nil {
		notFoundOrError(w, err, "precio", id)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ─────────────────────────────────────────────
// helper
// ─────────────────────────────────────────────

func notFoundOrError(w http.ResponseWriter, err error, recurso string, id int) {
	if errors.Is(err, services.ErrNoEncontrado) {
		services.NoEncontrado(w, recurso, id)
		return
	}
	services.ErrorMermoria(w, err, recurso, id)
}
