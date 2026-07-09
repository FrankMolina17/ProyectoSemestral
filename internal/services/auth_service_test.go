package services

import (
	"testing"

	"Sistem-Inte-Gestion-Control-Obras/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthService(t *testing.T) *AuthService {
	t.Helper()
	repo := storage.NuevoUsuarioStorage()
	return NuevoAuthService(repo)
}

func TestAuthService_Registrar(t *testing.T) {
	svc := setupAuthService(t)

	t.Run("valido -> success", func(t *testing.T) {
		u, err := svc.Registrar("test@test.com", "123456")
		require.NoError(t, err)
		assert.Equal(t, "test@test.com", u.Email)
		assert.Equal(t, "cliente", u.Rol)
	})

	t.Run("email duplicado -> error", func(t *testing.T) {
		_, err := svc.Registrar("test@test.com", "123456")
		assert.Error(t, err)
	})

	t.Run("email vacio -> error", func(t *testing.T) {
		_, err := svc.Registrar("", "123456")
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	svc := setupAuthService(t)

	_, err := svc.Registrar("user@test.com", "123456")
	require.NoError(t, err)

	t.Run("correcto -> token", func(t *testing.T) {
		token, err := svc.Login("user@test.com", "123456")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("password incorrecto -> error", func(t *testing.T) {
		_, err := svc.Login("user@test.com", "wrong")
		assert.Error(t, err)
	})

	t.Run("email no existe -> error", func(t *testing.T) {
		_, err := svc.Login("noexiste@test.com", "123456")
		assert.Error(t, err)
	})
}

func TestAuthService_RegistrarConRol(t *testing.T) {
	svc := setupAuthService(t)

	t.Run("admin -> success", func(t *testing.T) {
		u, err := svc.RegistrarConRol("admin@test.com", "123456", "admin")
		require.NoError(t, err)
		assert.Equal(t, "admin", u.Rol)
	})

	t.Run("rol invalido -> error", func(t *testing.T) {
		_, err := svc.RegistrarConRol("bad@test.com", "123456", "superadmin")
		assert.Error(t, err)
	})

	t.Run("email vacio -> error", func(t *testing.T) {
		_, err := svc.RegistrarConRol("", "123456", "admin")
		assert.Error(t, err)
	})
}

func TestAuthService_VerificarToken(t *testing.T) {
	svc := setupAuthService(t)

	_, err := svc.Registrar("vtest@test.com", "123456")
	require.NoError(t, err)

	t.Run("token valido", func(t *testing.T) {
		token, err := svc.Login("vtest@test.com", "123456")
		require.NoError(t, err)

		claims, err := svc.VerificarToken(token)
		require.NoError(t, err)
		assert.Equal(t, "cliente", claims.Rol)
		assert.Greater(t, claims.UsuarioID, 0)
	})

	t.Run("token invalido", func(t *testing.T) {
		_, err := svc.VerificarToken("token-invalido")
		assert.Error(t, err)
	})
}
