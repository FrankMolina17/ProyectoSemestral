package services

import (
	"errors"
	"strings"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretJwt = []byte("secreto")

// duracion del jwt
const DuracionJWT = time.Hour * 24

//Claims es el payload del JWT de un usuario.

type Claims struct {
	UsuarioID int    `json:"usuario_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type AutenticacionService struct {
	repo storage.UsuarioRepository
}

func NuevaAutenticacionService(repo storage.UsuarioRepository) *AutenticacionService {
	return &AutenticacionService{
		repo: repo,
	}
}

// registra un nuevo Usuario
func (s *AutenticacionService) RegistrarUsuario(email, password string) (*models.Usuario, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, ErrNombreVacio
	}
	if len(password) < 6 {
		return nil, errors.New("la contraseña debe tener al menos 6 caracteres")
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return nil, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.repo.CrearUsuario(models.EntradaUsuario{
		Email:    email,
		Password: string(hash),
	})
}

// Login de un usuario
func (s *AutenticacionService) Login(email, password string) (*models.Usuario, error) {
	usuario, ok := s.repo.BuscarUsuarioPorEmail(email)
	if !ok {
		return nil, ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(password)); err != nil {
		return nil, ErrCredencialesInvalidas
	}
	return &usuario, nil
}

func (s *AutenticacionService) ListarUsuarios() []*models.Usuario {
	return s.repo.ListarUsuarios()
}

func (s *AutenticacionService) ObtenerUsuarioPorID(id int) (*models.Usuario, bool) {
	return s.repo.ObtenerUsuarioPorID(id)
}

func (s *AutenticacionService) GenerarJWT(usuario models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: usuario.ID,
		Email:     usuario.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(DuracionJWT)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretJwt)
}

// validar Token
func (s *AutenticacionService) ValidarJWT(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretJwt, nil
	})
	if err != nil || !tok.Valid {
		return nil, ErrCredencialesInvalidas
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok {
		return nil, ErrCredencialesInvalidas
	}
	return claims, nil
}

var duracioToken = time.Hour * 24

const (
	secretoPorDefecto  = "incidencias-uleam-secreto-dev"
	duracionPorDefecto = 24 * time.Hour
)

type AuthService struct {
	repo     storage.UserRepository
	secreto  []byte
	duracion time.Duration
}

func NuevoAuthService(repo storage.UserRepository, opts ...AuthOption) *AuthService {
	s := &AuthService{
		repo:     repo,
		secreto:  []byte(secretoPorDefecto),
		duracion: duracionPorDefecto,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// AuthOption configura el AuthService
type AuthOption func(*AuthService)

func WithSecreto(secreto []byte) AuthOption {
	return func(a *AuthService) {
		if len(secreto) > 0 {
			a.secreto = secreto
		}
	}
}

func WithDuracionToken(d time.Duration) AuthOption {
	return func(a *AuthService) {
		if d > 0 {
			a.duracion = d
		}
	}
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
			CreatedAt:    time.Now(),
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
	return token.SignedString(secretJwt)
}

func (s *AuthService) VerificarToken(tokenSrt string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenSrt, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidos
		}
		return secretJwt, nil
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
