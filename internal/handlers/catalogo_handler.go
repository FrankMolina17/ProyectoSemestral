package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/go-chi/chi/v5"
)

type MaterialHandler struct {
	svc *services.MaterialService
}

func NuevoMaterialHandler(svc *services.MaterialService) *MaterialHandler {
	return &MaterialHandler{svc: svc}
}

func (h *MaterialHandler) ListandoMateriales(w http.ResponseWriter, r *http.Request) {
	materiales := h.svc.Listar()
	responderJSON(w, http.StatusOK, materiales)
}

func (h *MaterialHandler) ObtenerMaterialPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	m, err := h.svc.Obtener(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, m)
}

func (h *MaterialHandler) CreandoMaterial(w http.ResponseWriter, r *http.Request) {
	var m models.Material
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	resultado := h.svc.Crear(m)
	responderJSON(w, http.StatusCreated, resultado)
}

func (h *MaterialHandler) ActulizarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	var m models.Material
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	m.ID = id
	resultado, err := h.svc.Actualizar(m)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, resultado)
}

func (h *MaterialHandler) BorrarUnMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

type ManoObraHandler struct {
	svc *services.ManoObraService
}

func NuevoManoObraHandler(svc *services.ManoObraService) *ManoObraHandler {
	return &ManoObraHandler{svc: svc}
}

func (h *ManoObraHandler) ListarUnaManoObra(w http.ResponseWriter, r *http.Request) {
	lista := h.svc.Listar()
	responderJSON(w, http.StatusOK, lista)
}

func (h *ManoObraHandler) ObtenerUnaManoObraPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	m, err := h.svc.Obtener(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, m)
}

func (h *ManoObraHandler) CreandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	var m models.ManoObra
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	resultado := h.svc.Crear(m)
	responderJSON(w, http.StatusCreated, resultado)
}

func (h *ManoObraHandler) ActualizadoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	var m models.ManoObra
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	m.ID = id
	resultado, err := h.svc.Actualizar(m)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, resultado)
}

func (h *ManoObraHandler) BorrandoUnaManoObra(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

type EquipoHandler struct {
	svc *services.EquipoService
}

func NuevoEquipoHandler(svc *services.EquipoService) *EquipoHandler {
	return &EquipoHandler{svc: svc}
}

func (h *EquipoHandler) ListandoLosEquipos(w http.ResponseWriter, r *http.Request) {
	lista := h.svc.Listar()
	responderJSON(w, http.StatusOK, lista)
}

func (h *EquipoHandler) ObtenerUnEquipoPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	e, err := h.svc.Obtener(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, e)
}

func (h *EquipoHandler) CrearUnEquipo(w http.ResponseWriter, r *http.Request) {
	var e models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	resultado := h.svc.Crear(e)
	responderJSON(w, http.StatusCreated, resultado)
}

func (h *EquipoHandler) ActualizarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	var e models.Equipo
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	e.ID = id
	resultado, err := h.svc.Actualizar(e)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, resultado)
}

func (h *EquipoHandler) BorrarUnEquipo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

type PrecioHandler struct {
	svc *services.PreciosService
}

func NuevoPrecioHandler(svc *services.PreciosService) *PrecioHandler {
	return &PrecioHandler{svc: svc}
}

func (h *PrecioHandler) ListarDeLosPrecios(w http.ResponseWriter, r *http.Request) {
	lista := h.svc.Listar()
	responderJSON(w, http.StatusOK, lista)
}

func (h *PrecioHandler) CrearUnPrecio(w http.ResponseWriter, r *http.Request) {
	var p models.Precio
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	resultado := h.svc.Crear(p)
	responderJSON(w, http.StatusCreated, resultado)
}

func (h *PrecioHandler) PrecioVigentePorRecurso(w http.ResponseWriter, r *http.Request) {
	tipo := chi.URLParam(r, "tipo")
	recursoID, err := strconv.Atoi(chi.URLParam(r, "recursoID"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "recursoID inválido"})
		return
	}
	p, err := h.svc.PrecioVigente(tipo, recursoID)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func (h *PrecioHandler) HistorialPorRecurso(w http.ResponseWriter, r *http.Request) {
	tipo := chi.URLParam(r, "tipo")
	recursoID, err := strconv.Atoi(chi.URLParam(r, "recursoID"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "recursoID inválido"})
		return
	}
	historial := h.svc.HistorialPorRecurso(tipo, recursoID)
	responderJSON(w, http.StatusOK, historial)
}

func (h *PrecioHandler) ObtenerUnPrecioPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	p, err := h.svc.Obtener(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func (h *PrecioHandler) ActualizarUnPrecio(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	var p models.Precio
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	p.ID = id
	resultado, err := h.svc.Actualizar(p)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, resultado)
}

func (h *PrecioHandler) BorrarUnPrecio(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "id inválido"})
		return
	}
	if err := h.svc.Borrar(id); err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusNoContent, nil)
}

type ServerC struct {
	manoObraSvc *services.ManoObraService
	materialSvc *services.MaterialService
	equipoSvc   *services.EquipoService
	precioSvc   *services.PreciosService
	authSvc     *services.AutenticacionCatalogoService
}

func NuevoServerC(
	manoObraSvc *services.ManoObraService,
	materialSvc *services.MaterialService,
	equipoSvc *services.EquipoService,
	precioSvc *services.PreciosService,
	authSvc *services.AutenticacionCatalogoService,
) *ServerC {
	return &ServerC{
		manoObraSvc: manoObraSvc,
		materialSvc: materialSvc,
		equipoSvc:   equipoSvc,
		precioSvc:   precioSvc,
		authSvc:     authSvc,
	}
}

func (h *ServerC) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	token, err := h.authSvc.Login(creds.Email, creds.Password)
	if err != nil {
		responderJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *ServerC) RegistrarUser(w http.ResponseWriter, r *http.Request) {
	var u models.UsuarioCatalogo
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "cuerpo inválido"})
		return
	}
	responderJSON(w, http.StatusCreated, u)
}
