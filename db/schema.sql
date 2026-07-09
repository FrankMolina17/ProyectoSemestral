CREATE TABLE IF NOT EXISTS clientes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre      TEXT    NOT NULL,
    email       TEXT,
    telefono    TEXT,
    ruc         TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS proformas (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    obra_id         INTEGER NOT NULL,
    cliente_id      INTEGER REFERENCES clientes(id),
    nombre          TEXT    NOT NULL,
    estado          TEXT    NOT NULL DEFAULT 'borrador',
    pct_ganancia    REAL    NOT NULL DEFAULT 0,
    pct_imprevisto  REAL    NOT NULL DEFAULT 0,
    subtotal        REAL    NOT NULL DEFAULT 0,
    total           REAL    NOT NULL DEFAULT 0,
    creado_en       DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS proforma_items (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proforma_id     INTEGER NOT NULL REFERENCES proformas(id),
    tipo_recurso    TEXT    NOT NULL,
    recurso_id      INTEGER NOT NULL,
    descripcion     TEXT    NOT NULL,
    cantidad        REAL    NOT NULL,
    precio_promedio REAL    NOT NULL,
    subtotal        REAL    NOT NULL
);

CREATE TABLE IF NOT EXISTS nota_proformas (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    proforma_id INTEGER NOT NULL REFERENCES proformas(id),
    contenido   TEXT    NOT NULL,
    creado_en   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usuarios (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT    NOT NULL UNIQUE,
    password_hash TEXT    NOT NULL,
    creado_en     DATETIME DEFAULT CURRENT_TIMESTAMP
);
