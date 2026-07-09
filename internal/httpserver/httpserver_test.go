package httpserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNuevo_PorDefecto(t *testing.T) {
	s := Nuevo(http.DefaultServeMux)
	assert.Equal(t, ":8080", s.Addr)
	assert.Equal(t, 10*time.Second, s.ReadTimeout)
	assert.Equal(t, 10*time.Second, s.WriteTimeout)
	assert.Equal(t, 60*time.Second, s.IdleTimeout)
	assert.Equal(t, http.DefaultServeMux, s.Handler)
}

func TestConPuerto(t *testing.T) {
	t.Run("puerto valido", func(t *testing.T) {
		s := Nuevo(nil, ConPuerto(":3000"))
		assert.Equal(t, ":3000", s.Addr)
	})

	t.Run("puerto vacio no sobreescribe", func(t *testing.T) {
		s := Nuevo(nil, ConPuerto(""))
		assert.Equal(t, ":8080", s.Addr)
	})
}

func TestConReadTimeout(t *testing.T) {
	t.Run("timeout valido", func(t *testing.T) {
		s := Nuevo(nil, ConReadTimeout(5*time.Second))
		assert.Equal(t, 5*time.Second, s.ReadTimeout)
	})

	t.Run("timeout cero no sobreescribe", func(t *testing.T) {
		s := Nuevo(nil, ConReadTimeout(0))
		assert.Equal(t, 10*time.Second, s.ReadTimeout)
	})
}

func TestConWriteTimeout(t *testing.T) {
	t.Run("timeout valido", func(t *testing.T) {
		s := Nuevo(nil, ConWriteTimeout(5*time.Second))
		assert.Equal(t, 5*time.Second, s.WriteTimeout)
	})

	t.Run("timeout cero no sobreescribe", func(t *testing.T) {
		s := Nuevo(nil, ConWriteTimeout(0))
		assert.Equal(t, 10*time.Second, s.WriteTimeout)
	})
}

func TestMultiplesOpciones(t *testing.T) {
	s := Nuevo(
		http.DefaultServeMux,
		ConPuerto(":9090"),
		ConReadTimeout(3*time.Second),
		ConWriteTimeout(7*time.Second),
	)
	assert.Equal(t, ":9090", s.Addr)
	assert.Equal(t, 3*time.Second, s.ReadTimeout)
	assert.Equal(t, 7*time.Second, s.WriteTimeout)
}
