package storage

type Recursos struct {
	Almacen      *Storage
	Usuarios     *Storage
	BackendUsado string
	Cerrar       func() error
}

func Inicializar(rutaDB string) (*Recursos, error) {
	_ = rutaDB

	almacen := New()

	cerrar := func() error {
		return nil
	}

	return &Recursos{
		Almacen:      almacen,
		Usuarios:     almacen,
		BackendUsado: "memoria",
		Cerrar:       cerrar,
	}, nil
}
