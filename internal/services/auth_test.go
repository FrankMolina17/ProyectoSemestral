package services

import (
	"errors"
	"testing"
	

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	
	"github.com/stretchr/testify/assert"
)
// ─────────────────────────────────────────────
// USUARIOS (AutenticacionService)
// ─────────────────────────────────────────────

type mockUsuarioRepo struct {
	crearCalled bool
	porEmail    map[string]models.Usuario
}

func (m *mockUsuarioRepo) CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error) {
	m.crearCalled = true
	if m.porEmail == nil {
		m.porEmail = map[string]models.Usuario{}
	}
	u := models.Usuario{ID: len(m.porEmail) + 1, Email: in.Email, PasswordHash: in.Password}
	m.porEmail[in.Email] = u
	return &u, nil
}
func (m *mockUsuarioRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := m.porEmail[email]
	return u, ok
}
func (m *mockUsuarioRepo) ListarUsuarios() []*models.Usuario              { return nil }
func (m *mockUsuarioRepo) ObtenerUsuarioPorID(id int) (*models.Usuario, bool) { return nil, false }

func TestRegistrarUsuario_RechazaDatosInvalidos(t *testing.T) {
	casos := []struct {
		nombre   string
		email    string
		password string
	}{
		{"email vacio", "", "password123"},
		{"password corto", "user@example.com", "123"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			mock := &mockUsuarioRepo{}
			svc := NuevaAutenticacionService(mock, AuthOptions{})

			_, err := svc.RegistrarUsuario(tc.email, tc.password)

			assert.Error(t, err)
			assert.False(t, mock.crearCalled, "CrearUsuario no debe ser llamado cuando los datos son invalidos")
		})
	}
}

func TestRegistrarUsuario_DatoValido_LlamaAlRepositorio(t *testing.T) {
	mock := &mockUsuarioRepo{}
	svc := NuevaAutenticacionService(mock, AuthOptions{})

	_, err := svc.RegistrarUsuario("nuevo@example.com", "password123")

	assert.NoError(t, err)
	assert.True(t, mock.crearCalled, "CrearUsuario debe ser llamado cuando los datos son validos")
}

func TestRegistrarUsuario_EmailEnUso_RetornaError(t *testing.T) {
	mock := &mockUsuarioRepo{porEmail: map[string]models.Usuario{
		"repeat@example.com": {ID: 1, Email: "repeat@example.com", PasswordHash: "hash"},
	}}
	svc := NuevaAutenticacionService(mock, AuthOptions{})

	_, err := svc.RegistrarUsuario("repeat@example.com", "password123")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmailEnUso))
}

func TestLogin_Exito(t *testing.T) {
	mock := &mockUsuarioRepo{}
	svc := NuevaAutenticacionService(mock, AuthOptions{})

	_, err := svc.RegistrarUsuario("login@example.com", "password123")
	assert.NoError(t, err)

	u, err := svc.Login("login@example.com", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "login@example.com", u.Email)
}

func TestLogin_PasswordIncorrecto_RetornaError(t *testing.T) {
	mock := &mockUsuarioRepo{}
	svc := NuevaAutenticacionService(mock, AuthOptions{})

	_, _ = svc.RegistrarUsuario("login@example.com", "password123")

	_, err := svc.Login("login@example.com", "wrongpass")
	assert.Error(t, err)
}

func TestLogin_UsuarioNoExiste_RetornaError(t *testing.T) {
	mock := &mockUsuarioRepo{}
	svc := NuevaAutenticacionService(mock, AuthOptions{})

	_, err := svc.Login("noexiste@example.com", "password123")
	assert.Error(t, err)
}
