CREATE TABLE material (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    descripcion TEXT,
    unidad TEXT NOT NULL,
    precio_referencia TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE mano_obra (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    descripcion TEXT NOT NULL,
    categoria TEXT NOT NULL,
    unidad TEXT NOT NULL,
    costo_referencia TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE equipo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    tipo TEXT NOT NULL,
    unidad TEXT NOT NULL,
    costo_hora TEXT NOT NULL,
    disponible INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL
);

CREATE TABLE precios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recurso_tipo TEXT NOT NULL,
    recurso_id INTEGER NOT NULL,
    precio TEXT NOT NULL,
    fecha_vigencia TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE usuario (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL
);