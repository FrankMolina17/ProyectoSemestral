package services

import (
	"strings"
	"time"
	"errors"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretoJWTProforma = []byte("proformas-2026-secret")
var duracionTokenProforma = time.Hour * 24

type ClaimsProforma struct {
	UsuarioID int    `json:"uid"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo *storage.UsuarioStorage
}

func NuevoAuthService(repo *storage.UsuarioStorage) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrEmailVacio
	}

	if _, existe := s.repo.BuscarPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
		Rol:          "cliente", // rol por defecto; nunca se toma del request público
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, existe := s.repo.BuscarPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}

	return s.generarToken(u)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := ClaimsProforma{
		UsuarioID: u.ID,
		Rol:       u.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionTokenProforma)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWTProforma)
}

// VerificarToken ahora devuelve los claims completos (con Rol incluido)
// en lugar de solo el UsuarioID.
func (s *AuthService) VerificarToken(tokenStr string) (*ClaimsProforma, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimsProforma{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretoJWTProforma, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrCredencialesInvalidas
	}

	claims, ok := token.Claims.(*ClaimsProforma)
	if !ok {
		return nil, ErrCredencialesInvalidas
	}

	return claims, nil
}
func (s *AuthService) RegistrarConRol(email, password, rol string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrEmailVacio
	}

	if rol != "admin" && rol != "cliente" {
		return models.Usuario{}, errors.New("rol inválido: debe ser 'admin' o 'cliente'")
	}

	if _, existe := s.repo.BuscarPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}

	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
		Rol:          rol,
	})
}