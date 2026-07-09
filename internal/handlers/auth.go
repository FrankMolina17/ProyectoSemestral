package handlers

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/services"
	"encoding/json"
	"net/http"
)

type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *ServerC) RegistrarUser(w http.ResponseWriter, r *http.Request) {
	var credenciales credenciales
	if err := json.NewDecoder(r.Body).Decode(&credenciales); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos invalidos "+err.Error())
		return
	}
	usuario, err := s.Autenticacion.RegistrarUsuario(credenciales.Email, credenciales.Password)
	if err != nil {
		RespondError(w, StatusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, map[string]any{"id": usuario.ID, "email": usuario.Email})
}

func (s *ServerC) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var credenciales credenciales
	if err := json.NewDecoder(r.Body).Decode(&credenciales); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos invalidos "+err.Error())
		return
	}
	usuario, err := s.Autenticacion.Login(credenciales.Email, credenciales.Password)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := s.Autenticacion.GenerarJWT(*usuario)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *ServerC) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usuarios := s.Autenticacion.ListarUsuarios()
	RespondJSON(w, http.StatusOK, map[string]any{"data": usuarios})
}

func (s *ServerC) ObtenerUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, valid := services.ParaObtenerelID(w, r)
	if !valid {
		return
	}
	usuario, ok := s.Autenticacion.ObtenerUsuarioPorID(id)
	if !ok {
		services.NoEncontrado(w, "usuario", id)
		return
	}
	services.Ok(w, usuario)
}
