package storage

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsuarioStorage_CrearYAbrirEmail(t *testing.T) {
	s := NuevoUsuarioStorage()

	u, err := s.CrearUsuario(models.Usuario{
		Email:        "test@test.com",
		PasswordHash: "hash",
		Rol:          "cliente",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "test@test.com", u.Email)
	assert.Equal(t, "cliente", u.Rol)
	assert.False(t, u.CreatedAt.IsZero())

	_, err = s.CrearUsuario(models.Usuario{
		Email:        "test@test.com",
		PasswordHash: "otrohash",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ya está registrado")

	u2, ok := s.BuscarPorEmail("test@test.com")
	assert.True(t, ok)
	assert.Equal(t, "test@test.com", u2.Email)

	_, ok = s.BuscarPorEmail("noexiste@test.com")
	assert.False(t, ok)
}
