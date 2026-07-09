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
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

type AutenticacionService struct {
	repo       storage.UsuarioRepository
	secretJwt  []byte
	duracionJwt time.Duration
}

type AuthOptions struct {
	Secreto  []byte
	Duracion time.Duration
}

func NuevaAutenticacionService(repo storage.UsuarioRepository, opts AuthOptions) *AutenticacionService {
	if len(opts.Secreto) == 0 {
		opts.Secreto = []byte("secreto")
	}
	if opts.Duracion <= 0 {
		opts.Duracion = time.Hour * 24
	}
	return &AutenticacionService{
		repo:        repo,
		secretJwt:   opts.Secreto,
		duracionJwt: opts.Duracion,
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
		Rol:       usuario.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duracionJwt)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretJwt)
}

// validar Token
func (s *AutenticacionService) ValidarJWT(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return s.secretJwt, nil
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
