package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/repository"
	"Sistem-Inte-Gestion-Control-Obras/internal/services"

	"github.com/go-chi/chi/v5"
)

type ProformaHandler struct {
	svc *services.ProformaService
}

func NuevoHandler(svc *services.ProformaService) *ProformaHandler {
	return &ProformaHandler{svc: svc}
}

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

	creada, err := h.svc.CrearProforma(p)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) || errors.Is(err, services.ErrObraIDRequerido) {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

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

	p, err := h.svc.ObtenerPorID(id)
	if err != nil {
		if errors.Is(err, repository.ErrProformaNoEncontrada) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusOK, p)
}

func (h *ProformaHandler) ObtenerTodos(w http.ResponseWriter, r *http.Request) {
	lista, err := h.svc.ObtenerTodos()
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
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

	actualizada, err := h.svc.ActualizarProforma(id, datos)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, repository.ErrProformaNoEncontrada) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

	if err := h.svc.EliminarProforma(id); err != nil {
		if errors.Is(err, repository.ErrProformaNoEncontrada) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

	item.ProformaID = id
	creado, err := h.svc.AgregarItem(item)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDescripcionRequerida),
			errors.Is(err, services.ErrCantidadInvalida),
			errors.Is(err, services.ErrPrecioInvalido):
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, repository.ErrProformaNoEncontrada):
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
		default:
			responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
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

	items, err := h.svc.ObtenerItems(id)
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
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

	aprobada, err := h.svc.AprobarProforma(id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProformaYaAprobada):
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, repository.ErrProformaNoEncontrada):
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
		default:
			responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	responderJSON(w, http.StatusOK, aprobada)
}

func (h *ProformaHandler) CrearCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cuerpo del request inválido",
		})
		return
	}

	creado, err := h.svc.CrearCliente(c)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) || errors.Is(err, services.ErrRucRequerido) {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusCreated, creado)
}

func (h *ProformaHandler) ObtenerClientes(w http.ResponseWriter, r *http.Request) {
	lista, err := h.svc.ObtenerClientes()
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, lista)
}

func (h *ProformaHandler) ObtenerClientePorID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	c, err := h.svc.ObtenerClientePorID(id)
	if err != nil {
		if errors.Is(err, repository.ErrClienteNoEncontrado) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "cliente no encontrado"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusOK, c)
}

func (h *ProformaHandler) ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	var datos models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cuerpo del request inválido",
		})
		return
	}

	actualizado, err := h.svc.ActualizarCliente(id, datos)
	if err != nil {
		if errors.Is(err, services.ErrNombreRequerido) {
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, repository.ErrClienteNoEncontrado) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "cliente no encontrado"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusOK, actualizado)
}

func (h *ProformaHandler) EliminarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	if err := h.svc.EliminarCliente(id); err != nil {
		if errors.Is(err, repository.ErrClienteNoEncontrado) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "cliente no encontrado"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensaje": "cliente eliminado correctamente",
	})
}

func (h *ProformaHandler) ObtenerResumen(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	resumen, err := h.svc.ObtenerResumen(id)
	if err != nil {
		if errors.Is(err, repository.ErrProformaNoEncontrada) {
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
			return
		}
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responderJSON(w, http.StatusOK, resumen)
}

func (h *ProformaHandler) AgregarNota(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	var nota models.NotaProforma
	if err := json.NewDecoder(r.Body).Decode(&nota); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cuerpo del request inválido",
		})
		return
	}

	nota.ProformaID = id
	creada, err := h.svc.AgregarNota(nota)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrContenidoRequerido):
			responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, repository.ErrProformaNoEncontrada):
			responderJSON(w, http.StatusNotFound, map[string]string{"error": "proforma no encontrada"})
		default:
			responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	responderJSON(w, http.StatusCreated, creada)
}

func (h *ProformaHandler) ObtenerNotas(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "id inválido",
		})
		return
	}

	notas, err := h.svc.ObtenerNotas(id)
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, notas)
}
