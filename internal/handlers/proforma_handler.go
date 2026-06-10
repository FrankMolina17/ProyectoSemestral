package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/go-chi/chi/v5"
)

// ProformaHandler tiene acceso al storage
type ProformaHandler struct {
    storage *storage.ProformaStorage
}

// New crea el handler con su storage
func NuevoHandler(s *storage.ProformaStorage) *ProformaHandler {
    return &ProformaHandler{storage: s}
}

// helper para responder JSON — lo usarás en todos los handlers
func responderJSON(w http.ResponseWriter, status int, dato interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(dato)
}

func (h *ProformaHandler) CrearProforma(w http.ResponseWriter, r *http.Request) {
    var p models.Proforma

    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "cuerpo del request inválido",
        })
        return
    }

    // Validación básica
    if p.Nombre == "" {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "el campo nombre es requerido",
        })
        return
    }
    if p.ObraID == 0 {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "el campo obra_id es requerido",
        })
        return
    }

    creada := h.storage.CrearProforma(p)
    responderJSON(w, http.StatusCreated, creada)
}

func (h *ProformaHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    p, err := h.storage.ObtenerPorID(id)
    if err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    responderJSON(w, http.StatusOK, p)
}

func (h *ProformaHandler) ObtenerTodos(w http.ResponseWriter, r *http.Request) {
    lista := h.storage.ObtenerTodos()
    responderJSON(w, http.StatusOK, lista)
}

func (h *ProformaHandler) ActualizarProforma(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    var datos models.Proforma
    if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "cuerpo del request inválido",
        })
        return
    }

    if datos.Nombre == "" {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "el campo nombre es requerido",
        })
        return
    }

    actualizada, err := h.storage.ActualizarProforma(id, datos)
    if err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    responderJSON(w, http.StatusOK, actualizada)
}

func (h *ProformaHandler) EliminarProforma(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    if err := h.storage.EliminarProforma(id); err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    responderJSON(w, http.StatusOK, map[string]string{
        "mensaje": "proforma eliminada correctamente",
    })
}

func (h *ProformaHandler) AgregarItem(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    var item models.ProformaItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "cuerpo del request inválido",
        })
        return
    }

    // Validación básica
    if item.Descripcion == "" {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "el campo descripcion es requerido",
        })
        return
    }
    if item.Cantidad <= 0 {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "la cantidad debe ser mayor a 0",
        })
        return
    }
    if item.PrecioPromedio <= 0 {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "el precio debe ser mayor a 0",
        })
        return
    }

    item.ProformaID = id
    creado, err := h.storage.AgregarItem(item)
    if err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    responderJSON(w, http.StatusCreated, creado)
}

func (h *ProformaHandler) ObtenerItems(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    items := h.storage.ObtenerItems(id)
    responderJSON(w, http.StatusOK, items)
}

func (h *ProformaHandler) AprobarProforma(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "id inválido",
        })
        return
    }

    p, err := h.storage.ObtenerPorID(id)
    if err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    if p.Estado == "aprobada" {
        responderJSON(w, http.StatusBadRequest, map[string]string{
            "error": "la proforma ya está aprobada",
        })
        return
    }

    aprobada, err := h.storage.AprobarProforma(id)
    if err != nil {
        responderJSON(w, http.StatusNotFound, map[string]string{
            "error": "proforma no encontrada",
        })
        return
    }

    responderJSON(w, http.StatusOK, aprobada)
}