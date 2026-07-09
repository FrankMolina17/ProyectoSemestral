package storage

import (
	"errors"
	"sync"
	"time"

	"Sistem-Inte-Gestion-Control-Obras/internal/models"
)

type UsuarioStorage struct {
	mu       sync.Mutex
	usuarios map[string]models.Usuario
	nextID   int
}

func NuevoUsuarioStorage() *UsuarioStorage {
	return &UsuarioStorage{
		usuarios: make(map[string]models.Usuario),
		nextID:   1,
	}
}

func (s *UsuarioStorage) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, existe := s.usuarios[u.Email]; existe {
		return models.Usuario{}, errors.New("el email ya está registrado")
	}

	u.ID = s.nextID
	u.CreatedAt = time.Now()
	s.usuarios[u.Email] = u
	s.nextID++
	return u, nil
}

func (s *UsuarioStorage) BuscarPorEmail(email string) (models.Usuario, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, existe := s.usuarios[email]
	return u, existe
}
