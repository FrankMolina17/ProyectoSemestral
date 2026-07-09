package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCargar_ValoresPorDefecto(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("PUERTO")
		os.Unsetenv("PORT")
		os.Unsetenv("RUTA_DB")
		os.Unsetenv("JWT_SECRETO")
		os.Unsetenv("JWT_DURACION")
	})

	cfg := Cargar()
	assert.Equal(t, ":8080", cfg.Puerto)
	assert.Equal(t, "incidencia.db", cfg.RutaDB)
	assert.Equal(t, []byte("incidencias-uleam-secreto-dev-2026"), cfg.JWTSecreto)
	assert.Equal(t, 24*time.Hour, cfg.JWTDuracion)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.WriteTimeout)
}

func TestResolverPuerto(t *testing.T) {
	t.Run("PUERTO sin prefijo", func(t *testing.T) {
		os.Setenv("PUERTO", "3000")
		t.Cleanup(func() { os.Unsetenv("PUERTO") })
		assert.Equal(t, ":3000", resolverPuerto())
	})

	t.Run("PUERTO con prefijo", func(t *testing.T) {
		os.Setenv("PUERTO", ":4000")
		t.Cleanup(func() { os.Unsetenv("PUERTO") })
		assert.Equal(t, ":4000", resolverPuerto())
	})

	t.Run("PORT funciona como fallback", func(t *testing.T) {
		os.Setenv("PORT", "5000")
		t.Cleanup(func() { os.Unsetenv("PORT") })
		assert.Equal(t, ":5000", resolverPuerto())
	})

	t.Run("defecto 8080", func(t *testing.T) {
		assert.Equal(t, ":8080", resolverPuerto())
	})
}

func TestConDuracion(t *testing.T) {
	t.Run("variable valida", func(t *testing.T) {
		os.Setenv("TEST_DUR", "5m")
		t.Cleanup(func() { os.Unsetenv("TEST_DUR") })
		assert.Equal(t, 5*time.Minute, conDuracion("TEST_DUR", time.Hour))
	})

	t.Run("variable invalida -> defecto", func(t *testing.T) {
		os.Setenv("TEST_DUR", "no-valido")
		t.Cleanup(func() { os.Unsetenv("TEST_DUR") })
		assert.Equal(t, time.Hour, conDuracion("TEST_DUR", time.Hour))
	})

	t.Run("sin variable -> defecto", func(t *testing.T) {
		assert.Equal(t, time.Hour, conDuracion("NO_EXISTE", time.Hour))
	})
}

func TestConTexto(t *testing.T) {
	t.Run("variable presente", func(t *testing.T) {
		os.Setenv("TEST_TEXTO", "hola")
		t.Cleanup(func() { os.Unsetenv("TEST_TEXTO") })
		assert.Equal(t, "hola", conTexto("TEST_TEXTO", "mundo"))
	})

	t.Run("variable ausente -> defecto", func(t *testing.T) {
		assert.Equal(t, "mundo", conTexto("NO_EXISTE", "mundo"))
	})
}
