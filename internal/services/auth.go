package services

import (
	"Sistem-Inte-Gestion-Control-Obras/internal/models"
	"Sistem-Inte-Gestion-Control-Obras/internal/storage"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretoJWT = []byte("Lulita-2024")

var duracioToken = time.Hour * 24

type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}

}

func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrNombreVacio
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}
	return s.repo.CrearUsuario(
		models.Usuario{
			Email:        email,
			PasswordHash: string(hash),
			CreadoEn:     time.Now(),
		},
	)
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrEmailEnUso
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidos
	}
	return s.generarToken(u)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {

	Claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracioToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(secretoJWT)
}

func (s *AuthService) VerificarToken(tokenSrt string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenSrt, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidos
		}
		return secretoJWT, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidos
	}
	Claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidos
	}
	return Claims.UsuarioID, nil
}
