-- Crear cliente
INSERT INTO clientes (nombre, email, telefono, ruc)
VALUES ('Constructora del Pacífico', 'contacto@pacifico.com', '0991234567', '1234567890001');

-- Crear proforma
INSERT INTO proformas (obra_id, cliente_id, nombre, estado, pct_ganancia, pct_imprevisto, subtotal, total, creado_en)
VALUES (1, 1, 'Proforma Edificio Central', 'borrador', 0.10, 0.05, 0, 0, CURRENT_TIMESTAMP);

-- Agregar ítem material
INSERT INTO proforma_items (proforma_id, tipo_recurso, recurso_id, descripcion, cantidad, precio_promedio, subtotal)
VALUES (1, 'material', 1, 'Cemento Portland', 10, 12.50, 125.00);

-- Agregar ítem mano de obra
INSERT INTO proforma_items (proforma_id, tipo_recurso, recurso_id, descripcion, cantidad, precio_promedio, subtotal)
VALUES (1, 'mano_obra', 2, 'Albañil', 5, 30.00, 150.00);

-- Agregar ítem equipo
INSERT INTO proforma_items (proforma_id, tipo_recurso, recurso_id, descripcion, cantidad, precio_promedio, subtotal)
VALUES (1, 'equipo', 3, 'Mezcladora de concreto', 2, 50.00, 100.00);

-- Agregar nota
INSERT INTO nota_proformas (proforma_id, contenido, creado_en)
VALUES (1, 'Revisar precios antes de aprobar', CURRENT_TIMESTAMP);