CREATE TABLE IF NOT EXISTS clientes (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(255) NOT NULL,
    email       VARCHAR(255),
    telefono    VARCHAR(50),
    ruc         VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS proformas (
    id              SERIAL PRIMARY KEY,
    obra_id         INTEGER NOT NULL,
    cliente_id      INTEGER REFERENCES clientes(id),
    nombre          VARCHAR(255) NOT NULL,
    estado          VARCHAR(50) NOT NULL DEFAULT 'borrador',
    pct_ganancia    NUMERIC(5,2) NOT NULL DEFAULT 0,
    pct_imprevisto  NUMERIC(5,2) NOT NULL DEFAULT 0,
    subtotal        NUMERIC(12,2) NOT NULL DEFAULT 0,
    total           NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS proforma_items (
    id              SERIAL PRIMARY KEY,
    proforma_id     INTEGER NOT NULL REFERENCES proformas(id) ON DELETE CASCADE,
    tipo_recurso    VARCHAR(50) NOT NULL,
    recurso_id      INTEGER NOT NULL,
    descripcion     VARCHAR(255) NOT NULL,
    cantidad        NUMERIC(12,2) NOT NULL,
    precio_promedio NUMERIC(12,2) NOT NULL,
    subtotal        NUMERIC(12,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS nota_proformas (
    id          SERIAL PRIMARY KEY,
    proforma_id INTEGER NOT NULL REFERENCES proformas(id) ON DELETE CASCADE,
    contenido   TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS materiales (
    id       SERIAL PRIMARY KEY,
    nombre   VARCHAR(255) NOT NULL,
    unidad   VARCHAR(50) NOT NULL,
    categoria VARCHAR(50) NOT NULL DEFAULT 'material'
);

CREATE TABLE IF NOT EXISTS mano_obras (
    id       SERIAL PRIMARY KEY,
    nombre   VARCHAR(255) NOT NULL,
    unidad   VARCHAR(50) NOT NULL,
    categoria VARCHAR(50) NOT NULL DEFAULT 'mano_obra'
);

CREATE TABLE IF NOT EXISTS equipos (
    id       SERIAL PRIMARY KEY,
    nombre   VARCHAR(255) NOT NULL,
    unidad   VARCHAR(50) NOT NULL,
    categoria VARCHAR(50) NOT NULL DEFAULT 'equipo'
);

CREATE TABLE IF NOT EXISTS precio_recursos (
    id           SERIAL PRIMARY KEY,
    entidad_tipo VARCHAR(50) NOT NULL,
    entidad_id   INTEGER NOT NULL,
    valor        NUMERIC(12,2) NOT NULL,
    fecha_desde  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_hasta  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS obras (
    id         SERIAL PRIMARY KEY,
    nombre     VARCHAR(255) NOT NULL,
    direccion  TEXT,
    estado     VARCHAR(50) NOT NULL DEFAULT 'planificacion',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS incidencias (
    id           SERIAL PRIMARY KEY,
    titulo       VARCHAR(255) NOT NULL,
    descripcion  TEXT,
    entidad_tipo VARCHAR(50) NOT NULL,
    entidad_id   INTEGER NOT NULL,
    prioridad    VARCHAR(50) NOT NULL DEFAULT 'media',
    estado       VARCHAR(50) NOT NULL DEFAULT 'abierta',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS usuarios (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    rol           VARCHAR(50) NOT NULL DEFAULT 'cliente',
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
