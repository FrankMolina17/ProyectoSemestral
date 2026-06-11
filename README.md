# ProyectoSemestral

Estructura del Sistema. 

M1 — Catálogo de Recursos. 

Gestiona: materiales, mano de obra, equipos, historial de precios.  

Permite mantener recursos organizados y actualizados para los cálculos de presupuestos. 

# Módulo 2 — Proformas y Cálculo

**Proyecto:** Sistema de Gestión y Control de Obras  
**Curso:** TDI-601 Aplicaciones Web II · 2026-1 · Paralelo B  
**Integrante:** Franklin Vicente Molina Ponce  
**Puerto:** `8082`

---

## ¿Qué hace este módulo?

Gestiona los presupuestos (proformas) de obras de construcción. Permite crear proformas, agregar recursos (materiales, mano de obra, equipos) y calcular el costo total usando el **método de precio promedio**.

El total se calcula automáticamente con la fórmula:

```
subtotal   = Σ (cantidad × precio_promedio) de cada ítem
ganancia   = subtotal × pct_ganancia
imprevisto = subtotal × pct_imprevisto
total      = subtotal + ganancia + imprevisto
```

---

## Estructura del proyecto

```
ProyectoSemestral-Modulo2-Proformas/
├── cmd/
│   └── api/
│       └── main.go                  ← punto de entrada, rutas Chi
├── internal/
│   ├── models/
│   │   └── proforma.go              ← structs Proforma y ProformaItem
│   ├── handlers/
│   │   └── proforma_handler.go      ← 9 handlers HTTP
│   └── storage/
│       └── proforma_storage.go      ← almacenamiento en memoria
├── go.mod
├── go.sum
└── README.md
```

---

## Cómo levantar el servidor

```bash
# Desde la carpeta del módulo
go run cmd/api/main.go

# O con puerto específico
PORT=8082 go run cmd/api/main.go
```

El servidor responde en: `http://localhost:8082`

Verificar que está corriendo:
```bash
curl http://localhost:8082/health
```

---

## Endpoints disponibles

| Método | Ruta | Descripción | Status |
|--------|------|-------------|--------|
| GET | `/health` | Verificar que el servidor está activo | 200 |
| POST | `/api/v1/proformas` | Crear una proforma nueva | 201 |
| GET | `/api/v1/proformas` | Listar todas las proformas | 200 |
| GET | `/api/v1/proformas/{id}` | Obtener una proforma por ID | 200 / 404 |
| PUT | `/api/v1/proformas/{id}` | Actualizar nombre y porcentajes | 200 / 404 |
| DELETE | `/api/v1/proformas/{id}` | Eliminar una proforma | 200 / 404 |
| POST | `/api/v1/proformas/{id}/items` | Agregar ítem y recalcular total | 201 / 400 / 404 |
| GET | `/api/v1/proformas/{id}/items` | Listar ítems de una proforma | 200 |
| GET | `/api/v1/proformas/{id}/resumen` | Ver desglose del costo total | 200 / 404 |
| PUT | `/api/v1/proformas/{id}/aprobar` | Cambiar estado a aprobada | 200 / 404 |

---

## Datos de prueba

### 1. Health check
```bash
curl http://localhost:8082/health
```

**Respuesta esperada `200 OK`:**
```json
{"estado":"ok","modulo":"proformas"}
```

---

### 2. Crear proforma
```bash
curl -X POST http://localhost:8082/api/v1/proformas \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Proforma Edificio Central",
    "obra_id": 1,
    "pct_ganancia": 0.10,
    "pct_imprevisto": 0.05
  }'
```

**Respuesta esperada `201 Created`:**
```json
{
  "id": 1,
  "obra_id": 1,
  "nombre": "Proforma Edificio Central",
  "estado": "borrador",
  "pct_ganancia": 0.1,
  "pct_imprevisto": 0.05,
  "subtotal": 0,
  "total": 0,
  "creado_en": "2026-06-10T..."
}
```

---

### 3. Crear segunda proforma
```bash
curl -X POST http://localhost:8082/api/v1/proformas \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Proforma Puente Rio Verde",
    "obra_id": 2,
    "pct_ganancia": 0.15,
    "pct_imprevisto": 0.08
  }'
```

---

### 4. Listar todas las proformas
```bash
curl http://localhost:8082/api/v1/proformas
```

**Respuesta esperada `200 OK`:**
```json
[
  { "id": 1, "nombre": "Proforma Edificio Central", "estado": "borrador", ... },
  { "id": 2, "nombre": "Proforma Puente Rio Verde", "estado": "borrador", ... }
]
```

---

### 5. Obtener proforma por ID
```bash
curl http://localhost:8082/api/v1/proformas/1
```

**Respuesta esperada `200 OK`:** datos de la proforma 1.

---

### 6. Agregar ítem — material
```bash
curl -X POST http://localhost:8082/api/v1/proformas/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "descripcion": "Cemento Portland",
    "tipo_recurso": "material",
    "recurso_id": 1,
    "cantidad": 10,
    "precio_promedio": 12.50
  }'
```

**Respuesta esperada `201 Created`:**
```json
{
  "id": 1,
  "proforma_id": 1,
  "tipo_recurso": "material",
  "descripcion": "Cemento Portland",
  "cantidad": 10,
  "precio_promedio": 12.5,
  "subtotal": 125
}
```

---

### 7. Agregar ítem — mano de obra
```bash
curl -X POST http://localhost:8082/api/v1/proformas/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "descripcion": "Albañil",
    "tipo_recurso": "mano_obra",
    "recurso_id": 2,
    "cantidad": 5,
    "precio_promedio": 30.00
  }'
```

**Subtotal ítem:** `5 × 30 = 150`

---

### 8. Agregar ítem — equipo
```bash
curl -X POST http://localhost:8082/api/v1/proformas/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "descripcion": "Mezcladora de concreto",
    "tipo_recurso": "equipo",
    "recurso_id": 3,
    "cantidad": 2,
    "precio_promedio": 50.00
  }'
```

**Subtotal ítem:** `2 × 50 = 100`

---

### 9. Ver ítems de la proforma
```bash
curl http://localhost:8082/api/v1/proformas/1/items
```

**Respuesta esperada `200 OK`:** lista de los 3 ítems agregados.

---

### 10. Ver resumen con totales calculados
```bash
curl http://localhost:8082/api/v1/proformas/1/resumen
```

**Respuesta esperada `200 OK`:**
```json
{
  "proforma_id": 1,
  "nombre": "Proforma Edificio Central",
  "estado": "borrador",
  "subtotal": 375,
  "ganancia": 37.5,
  "imprevisto": 18.75,
  "total": 431.25,
  "pct_ganancia": 0.1,
  "pct_imprevisto": 0.05
}
```

**Verificación del cálculo:**
```
subtotal   = 125 + 150 + 100 = 375.00
ganancia   = 375 × 0.10     =  37.50
imprevisto = 375 × 0.05     =  18.75
total      = 375 + 37.50 + 18.75 = 431.25 ✓
```

---

### 11. Actualizar proforma
```bash
curl -X PUT http://localhost:8082/api/v1/proformas/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Proforma Edificio Central v2",
    "pct_ganancia": 0.12,
    "pct_imprevisto": 0.05
  }'
```

**Respuesta esperada `200 OK`:** proforma con datos actualizados.

---

### 12. Aprobar proforma
```bash
curl -X PUT http://localhost:8082/api/v1/proformas/1/aprobar
```

**Respuesta esperada `200 OK`:**
```json
{
  "id": 1,
  "estado": "aprobada",
  ...
}
```

---

### 13. Eliminar proforma
```bash
curl -X DELETE http://localhost:8082/api/v1/proformas/2
```

**Respuesta esperada `200 OK`:**
```json
{"mensaje": "proforma eliminada correctamente"}
```

---

## Casos de error

### ID inexistente → 404
```bash
curl http://localhost:8082/api/v1/proformas/999
```
```json
{"error": "proforma no encontrada"}
```

### Sin campo nombre → 400
```bash
curl -X POST http://localhost:8082/api/v1/proformas \
  -H "Content-Type: application/json" \
  -d '{"obra_id": 1}'
```
```json
{"error": "el campo nombre es requerido"}
```

### Sin obra_id → 400
```bash
curl -X POST http://localhost:8082/api/v1/proformas \
  -H "Content-Type: application/json" \
  -d '{"nombre": "Proforma sin obra"}'
```
```json
{"error": "el campo obra_id es requerido"}
```

### Tipo de recurso inválido → 400
```bash
curl -X POST http://localhost:8082/api/v1/proformas/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "descripcion": "Algo",
    "tipo_recurso": "invalido",
    "cantidad": 5,
    "precio_promedio": 10
  }'
```
```json
{"error": "tipo_recurso debe ser material, mano_obra o equipo"}
```

### Cantidad negativa → 400
```bash
curl -X POST http://localhost:8082/api/v1/proformas/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "descripcion": "Cemento",
    "tipo_recurso": "material",
    "cantidad": -5,
    "precio_promedio": 12.50
  }'
```
```json
{"error": "la cantidad debe ser mayor a 0"}
```

### ID inválido (texto en lugar de número) → 400
```bash
curl http://localhost:8082/api/v1/proformas/abc
```
```json
{"error": "id inválido"}
```

---

## Método de costeo — Precio Promedio

El módulo implementa el método de **precio promedio** para calcular el costo de una proforma.

Cada ítem tiene un `precio_promedio` que representa el precio histórico promedio del recurso. El sistema calcula automáticamente:

1. **Subtotal del ítem:** `cantidad × precio_promedio`
2. **Subtotal de la proforma:** suma de todos los subtotales de ítems
3. **Ganancia:** `subtotal × pct_ganancia`
4. **Imprevisto:** `subtotal × pct_imprevisto`  
5. **Total final:** `subtotal + ganancia + imprevisto`

El recálculo es automático cada vez que se agrega un ítem — el cliente no necesita calcular nada.

---

## Tecnologías usadas

- **Go** — lenguaje de programación
- **Chi** — router HTTP (`github.com/go-chi/chi/v5`)
- **Almacenamiento en memoria** — mapa Go con `sync.Mutex`
M3 — Obras e Incidencias. 

Gestiona: obras, seguimiento de ejecución, incidencias y problemas registrados durante la construcción. Permite mejorar el control y trazabilidad de la obra. 

Users/Auth. 

Módulo compartido encargado de: autenticación, roles, seguridad mediante JWT.  

 


