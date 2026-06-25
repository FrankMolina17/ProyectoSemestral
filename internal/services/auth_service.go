package services

import (
	"strings"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretoJWT = []byte("proformas-2026-secret")
var duracionToken = time.Hour * 24

type Claims struct {
	UsuarioID int `json:"uid"`
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
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, existe := s.repo.BuscarPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidos
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidos
	}

	return s.generarToken(u)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWT)
}

func (s *AuthService) VerificarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidos
		}
		return secretoJWT, nil
	})

	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidos
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidos
	}

	return claims.UsuarioID, nil
}
