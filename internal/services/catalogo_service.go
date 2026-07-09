package services

import (
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

type MaterialService struct {
	store *storage.CatalogoStorage
}

func NuevoMaterialService(store *storage.CatalogoStorage) *MaterialService {
	return &MaterialService{store: store}
}

func (s *MaterialService) Listar() []models.Material {
	return s.store.ListarMateriales()
}

func (s *MaterialService) Obtener(id int) (models.Material, error) {
	m, ok := s.store.ObtenerMaterial(id)
	if !ok {
		return models.Material{}, ErrRecursoNoEncontrado
	}
	return m, nil
}

func (s *MaterialService) Crear(m models.Material) models.Material {
	return s.store.CrearMaterial(m)
}

func (s *MaterialService) Actualizar(m models.Material) (models.Material, error) {
	actualizado, ok := s.store.ActualizarMaterial(m)
	if !ok {
		return models.Material{}, ErrRecursoNoEncontrado
	}
	return actualizado, nil
}

func (s *MaterialService) Borrar(id int) error {
	if !s.store.BorrarMaterial(id) {
		return ErrRecursoNoEncontrado
	}
	return nil
}

type ManoObraService struct {
	store *storage.CatalogoStorage
}

func NuevoManoObraService(store *storage.CatalogoStorage) *ManoObraService {
	return &ManoObraService{store: store}
}

func (s *ManoObraService) Listar() []models.ManoObra {
	return s.store.ListarManoObras()
}

func (s *ManoObraService) Obtener(id int) (models.ManoObra, error) {
	m, ok := s.store.ObtenerManoObra(id)
	if !ok {
		return models.ManoObra{}, ErrRecursoNoEncontrado
	}
	return m, nil
}

func (s *ManoObraService) Crear(m models.ManoObra) models.ManoObra {
	return s.store.CrearManoObra(m)
}

func (s *ManoObraService) Actualizar(m models.ManoObra) (models.ManoObra, error) {
	actualizado, ok := s.store.ActualizarManoObra(m)
	if !ok {
		return models.ManoObra{}, ErrRecursoNoEncontrado
	}
	return actualizado, nil
}

func (s *ManoObraService) Borrar(id int) error {
	if !s.store.BorrarManoObra(id) {
		return ErrRecursoNoEncontrado
	}
	return nil
}

type EquipoService struct {
	store *storage.CatalogoStorage
}

func NuevoEquipoService(store *storage.CatalogoStorage) *EquipoService {
	return &EquipoService{store: store}
}

func (s *EquipoService) Listar() []models.Equipo {
	return s.store.ListarEquipos()
}

func (s *EquipoService) Obtener(id int) (models.Equipo, error) {
	e, ok := s.store.ObtenerEquipo(id)
	if !ok {
		return models.Equipo{}, ErrRecursoNoEncontrado
	}
	return e, nil
}

func (s *EquipoService) Crear(e models.Equipo) models.Equipo {
	return s.store.CrearEquipo(e)
}

func (s *EquipoService) Actualizar(e models.Equipo) (models.Equipo, error) {
	actualizado, ok := s.store.ActualizarEquipo(e)
	if !ok {
		return models.Equipo{}, ErrRecursoNoEncontrado
	}
	return actualizado, nil
}

func (s *EquipoService) Borrar(id int) error {
	if !s.store.BorrarEquipo(id) {
		return ErrRecursoNoEncontrado
	}
	return nil
}

type PreciosService struct {
	store *storage.CatalogoStorage
}

func NuevoPreciosService(store *storage.CatalogoStorage) *PreciosService {
	return &PreciosService{store: store}
}

func (s *PreciosService) Listar() []models.Precio {
	return s.store.ListarPrecios()
}

func (s *PreciosService) Obtener(id int) (models.Precio, error) {
	p, ok := s.store.ObtenerPrecio(id)
	if !ok {
		return models.Precio{}, ErrRecursoNoEncontrado
	}
	return p, nil
}

func (s *PreciosService) Crear(p models.Precio) models.Precio {
	return s.store.CrearPrecio(p)
}

func (s *PreciosService) Actualizar(p models.Precio) (models.Precio, error) {
	actualizado, ok := s.store.ActualizarPrecio(p)
	if !ok {
		return models.Precio{}, ErrRecursoNoEncontrado
	}
	return actualizado, nil
}

func (s *PreciosService) Borrar(id int) error {
	if !s.store.BorrarPrecio(id) {
		return ErrRecursoNoEncontrado
	}
	return nil
}

func (s *PreciosService) PrecioVigente(tipo string, recursoID int) (models.Precio, error) {
	p, ok := s.store.PrecioVigente(tipo, recursoID)
	if !ok {
		return models.Precio{}, ErrRecursoNoEncontrado
	}
	return p, nil
}

func (s *PreciosService) HistorialPorRecurso(tipo string, recursoID int) []models.Precio {
	return s.store.PreciosPorRecurso(tipo, recursoID)
}

var secretoJWTCatalogo = []byte("catalogo-2026-secret")
var duracionTokenCatalogo = time.Hour * 24

type ClaimsCatalogo struct {
	UsuarioID int    `json:"uid"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type AutenticacionCatalogoService struct {
	store *storage.CatalogoStorage
}

func NuevaAutenticacionService(store *storage.CatalogoStorage) *AutenticacionCatalogoService {
	return &AutenticacionCatalogoService{store: store}
}

func (s *AutenticacionCatalogoService) Login(email, password string) (string, error) {
	user, ok := s.store.BuscarUsuarioCatalogo(email, password)
	if !ok {
		return "", ErrCredencialesInvalidos
	}
	return s.generarToken(user)
}

func (s *AutenticacionCatalogoService) generarToken(u models.UsuarioCatalogo) (string, error) {
	claims := ClaimsCatalogo{
		UsuarioID: u.ID,
		Email:     u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionTokenCatalogo)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWTCatalogo)
}

func (s *AutenticacionCatalogoService) VerificarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsCatalogo{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidos
		}
		return secretoJWTCatalogo, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidos
	}
	claims, ok := token.Claims.(*ClaimsCatalogo)
	if !ok {
		return 0, ErrCredencialesInvalidos
	}
	return claims.UsuarioID, nil
}
