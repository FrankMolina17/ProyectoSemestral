package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuración
type Config struct {
	Puerto       string
	RutaDB       string
	JWTSecreto   []byte
	JWTDuracion  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Cargar() Config {
	_ = godotenv.Load() // Ignora error si no existe .env

	return Config{
		Puerto:       conTexto("PUERTO", ":8080"),
		RutaDB:       conTexto("RUTA_DB", "incidencia.db"),
		JWTSecreto:   []byte(conTexto("JWT_SECRETO", "incidencias-uleam-secreto-dev-2026")),
		JWTDuracion:  conDuracion("JWT_DURACION", 24*time.Hour),
		ReadTimeout:  conDuracion("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: conDuracion("HTTP_WRITE_TIMEOUT", 10*time.Second),
	}
}

func conTexto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

func conDuracion(clave string, porDefecto time.Duration) time.Duration {
	v := os.Getenv(clave)
	if v == "" {
		return porDefecto
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return porDefecto
	}
	return d
<<<<<<< HEAD
}
=======
}
>>>>>>> Modulo1/Catalogo
