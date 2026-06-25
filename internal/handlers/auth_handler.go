package handlers

import (
	"encoding/json"
	"net/http"

	"Sistem-Inte-Gestion-Control-Obras/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NuevoAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Registrar(w http.ResponseWriter, r *http.Request) {
	var creds credenciales

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cuerpo del request inválido",
		})
		return
	}

	u, err := h.service.Registrar(creds.Email, creds.Password)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	responderJSON(w, http.StatusCreated, map[string]interface{}{
		"id":        u.ID,
		"email":     u.Email,
		"creado_en": u.CreadoEn,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds credenciales

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{
			"error": "cuerpo del request inválido",
		})
		return
	}

	token, err := h.service.Login(creds.Email, creds.Password)
	if err != nil {
		responderJSON(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
