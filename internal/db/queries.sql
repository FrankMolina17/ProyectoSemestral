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
