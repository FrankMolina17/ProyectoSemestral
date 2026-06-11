# ProyectoSemestral

Estructura del Sistema. 

M1 — Catálogo de Recursos. 

Gestiona: materiales, mano de obra, equipos, historial de precios.  

Permite mantener recursos organizados y actualizados para los cálculos de presupuestos. 

M2 — Proformas y Cálculo. 
Gestiona: presupuestos, cálculos automáticos, costos de materiales, mano de obra y equipos.  
Este módulo permite generar proformas y estimar costos reales de una obra. 
Endpoints disponibles

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | /api/v1/proformas | Crear proforma |
| GET | /api/v1/proformas | Listar proformas |
| GET | /api/v1/proformas/{id} | Obtener proforma por ID |
| PUT | /api/v1/proformas/{id} | Actualizar proforma |
| DELETE | /api/v1/proformas/{id} | Eliminar proforma |
| POST | /api/v1/proformas/{id}/items | Agregar ítem |
| GET | /api/v1/proformas/{id}/items | Listar ítems |
| PUT | /api/v1/proformas/{id}/aprobar | Aprobar proforma |

M2 - Método de costeo
El módulo usa precio promedio. El total se calcula como:
subtotal + (subtotal × pct_ganancia) + (subtotal × pct_imprevisto)

M3 — Obras e Incidencias. 

Gestiona: obras, seguimiento de ejecución, incidencias y problemas registrados durante la construcción. Permite mejorar el control y trazabilidad de la obra. 

Users/Auth. 

Módulo compartido encargado de: autenticación, roles, seguridad mediante JWT.  

 


