package services

import (
	"testing"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	usuarios []models.Usuario
	nextID   int
}

func (m *mockUserRepo) CrearUsuario(in models.EntradaUsuario) (*models.Usuario, error) {
	m.nextID++
	u := &models.Usuario{
		ID:           m.nextID,
		Email:        in.Email,
		PasswordHash: in.Password,
		Rol:          "cliente",
	}
	m.usuarios = append(m.usuarios, *u)
	return u, nil
}

func (m *mockUserRepo) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	for _, u := range m.usuarios {
		if u.Email == email {
			return u, true
		}
	}
	return models.Usuario{}, false
}

func (m *mockUserRepo) ListarUsuarios() []*models.Usuario {
	res := make([]*models.Usuario, len(m.usuarios))
	for i := range m.usuarios {
		res[i] = &m.usuarios[i]
	}
	return res
}

func (m *mockUserRepo) ObtenerUsuarioPorID(id int) (*models.Usuario, bool) {
	for _, u := range m.usuarios {
		if u.ID == id {
			return &u, true
		}
	}
	return nil, false
}

func TestAutenticacionService_ListarYBuscar(t *testing.T) {
	svc := NuevaAutenticacionService(&mockUserRepo{}, AuthOptions{Secreto: []byte("test-secret"), Duracion: time.Hour})

	svc.RegistrarUsuario("a@test.com", "123456")
	svc.RegistrarUsuario("b@test.com", "123456")

	t.Run("ListarUsuarios", func(t *testing.T) {
		usuarios := svc.ListarUsuarios()
		assert.Len(t, usuarios, 2)
	})

	t.Run("ObtenerUsuarioPorID existente", func(t *testing.T) {
		u, ok := svc.ObtenerUsuarioPorID(1)
		require.True(t, ok)
		assert.Equal(t, "a@test.com", u.Email)
	})

	t.Run("ObtenerUsuarioPorID inexistente", func(t *testing.T) {
		_, ok := svc.ObtenerUsuarioPorID(999)
		assert.False(t, ok)
	})
}

func TestAutenticacionService_GenerarJWT(t *testing.T) {
	svc := NuevaAutenticacionService(&mockUserRepo{}, AuthOptions{Secreto: []byte("test-secret"), Duracion: time.Hour})

	usuario := models.Usuario{ID: 1, Email: "test@test.com", Rol: "admin"}
	token, err := svc.GenerarJWT(usuario)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAutenticacionService_ValidarJWT(t *testing.T) {
	svc := NuevaAutenticacionService(&mockUserRepo{}, AuthOptions{Secreto: []byte("test-secret"), Duracion: time.Hour})

	usuario := models.Usuario{ID: 1, Email: "test@test.com", Rol: "admin"}
	token, err := svc.GenerarJWT(usuario)
	require.NoError(t, err)

	t.Run("valido", func(t *testing.T) {
		claims, err := svc.ValidarJWT(token)
		require.NoError(t, err)
		assert.Equal(t, 1, claims.UsuarioID)
		assert.Equal(t, "test@test.com", claims.Email)
		assert.Equal(t, "admin", claims.Rol)
	})

	t.Run("token invalido", func(t *testing.T) {
		_, err := svc.ValidarJWT("token-invalido")
		assert.Error(t, err)
	})

	t.Run("firma incorrecta", func(t *testing.T) {
		svc2 := NuevaAutenticacionService(&mockUserRepo{}, AuthOptions{Secreto: []byte("otro-secreto"), Duracion: time.Hour})
		_, err := svc2.ValidarJWT(token)
		assert.Error(t, err)
	})
}
