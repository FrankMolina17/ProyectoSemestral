CREATE TABLE material (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    descripcion TEXT,
    unidad TEXT NOT NULL,
    precio_referencia TEXT NOT NULL,
    created_at TEXT NOT NULL
);
