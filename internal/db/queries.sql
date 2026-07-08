-- name: ListarMateriales :many
SELECT * FROM material ORDER BY id;

-- name: BuscarMaterialPorID :one
SELECT * FROM material WHERE id = ?;

-- name: CrearMateriales :one
INSERT INTO material (nombre, descripcion, unidad, precio_referencia, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ActualizarMateriales :one
UPDATE material SET nombre = ?, descripcion = ?, unidad = ?, precio_referencia = ?
WHERE id = ?
RETURNING *;

-- name: EliminarMateriales :execrows
DELETE FROM material WHERE id = ?;

-- name: ListarManoObra :many
SELECT * FROM mano_obra ORDER BY id;

-- name: BuscarManoObraPorID :one
SELECT * FROM mano_obra WHERE id = ?;

-- name: CrearManoObra :one
INSERT INTO mano_obra (descripcion, categoria, unidad, costo_referencia, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ActualizarManoObra :one
UPDATE mano_obra SET descripcion = ?, categoria = ?, unidad = ?, costo_referencia = ?
WHERE id = ?
RETURNING *;

-- name: EliminarManoObra :execrows
DELETE FROM mano_obra WHERE id = ?;

-- name: ListarEquipos :many
SELECT * FROM equipo ORDER BY id;

-- name: BuscarEquipoPorID :one
SELECT * FROM equipo WHERE id = ?;

-- name: CrearEquipo :one
INSERT INTO equipo (nombre, tipo, unidad, costo_hora, disponible, created_at)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ActualizarEquipo :one
UPDATE equipo SET nombre = ?, tipo = ?, unidad = ?, costo_hora = ?, disponible = ?
WHERE id = ?
RETURNING *;

-- name: EliminarEquipo :execrows
DELETE FROM equipo WHERE id = ?;

-- name: ListarPrecios :many
SELECT * FROM precios ORDER BY id;

-- name: BuscarPrecioPorID :one
SELECT * FROM precios WHERE id = ?;

-- name: ListarPreciosPorRecurso :many
SELECT * FROM precios WHERE recurso_tipo = ? AND recurso_id = ? ORDER BY fecha_vigencia DESC;

-- name: CrearPrecio :one
INSERT INTO precios (recurso_tipo, recurso_id, precio, fecha_vigencia, created_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: ActualizarPrecio :one
UPDATE precios SET recurso_tipo = ?, recurso_id = ?, precio = ?, fecha_vigencia = ?
WHERE id = ?
RETURNING *;

-- name: EliminarPrecio :execrows
DELETE FROM precios WHERE id = ?;

-- name: ListarUsuarios :many
SELECT * FROM usuario ORDER BY id;

-- name: BuscarUsuarioPorID :one
SELECT * FROM usuario WHERE id = ?;

-- name: BuscarUsuarioPorEmail :one
SELECT * FROM usuario WHERE email = ?;

-- name: CrearUsuario :one
INSERT INTO usuario (email, password_hash, created_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: ActualizarUsuario :one
UPDATE usuario SET email = ?, password_hash = ?
WHERE id = ?
RETURNING *;

-- name: EliminarUsuario :execrows
DELETE FROM usuario WHERE id = ?;