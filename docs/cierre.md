# Documento de Cierre — Sistema de Gestión y Control de Obras

## Integrantes
- Franklin Molina (M2 — Proformas)
- Melani Molina (M1 — Catálogo)
- Carlos Bailón (M3 — Obras e Incidencias)

## ¿Qué implementamos?
API RESTful en Go que integra tres módulos: catálogo de recursos, proformas con cálculo de costos y obras/incidencias. Cada módulo expone su propio conjunto de endpoints bajo `api/v1/`. La autenticación es centralizada con JWT y dos roles (admin, cliente).

## Principales dificultades
1. **Integración de módulos** — Al inicio cada integrante trabajó separado; unificar routers, estructuras y nomenclatura fue el desafío más grande.
2. **Polimorfismo en GORM** — Las incidencias y precios pueden pertenecer a múltiples entidades; implementar `EntidadTipo` + `EntidadID` en lugar de claves foráneas directas requirió varios intentos.
3. **Cobertura de tests** — Pasar del 11% al 50%+ implicó escribir tests con mock de servicios y probar casos borde como IDs inválidos y roles no autorizados.
4. **Convenciones Go** — Aprendizaje del manejo de contextos, middleware Chi, y evitar panic con `LogError` en handlers.

## Cambios respecto a HITO 1 y 2
- Se migró de almacenamiento en memoria a GORM (SQLite/PostgreSQL).
- Se agregó autenticación JWT y protección por roles.
- Se unificaron las bases de código de tres repos separados en uno solo con `cmd/api/main.go` centralizado.
- Se agregó Docker multi-stage y CI/CD con GitHub Actions.

## ¿Qué haríamos diferente?
- Usar un router con generación de OpenAPI (Huma v2) desde el inicio para documentación automática.
- Definir las interfaces de storage ANTES de escribir handlers.
- Hacer code reviews con PRs desde la primera semana.

## Próximos pasos
- Frontend web (React/Vue) para consumir la API.
- WebSockets para notificaciones en tiempo real al registrar incidencias.
- Reportes PDF exportables desde las proformas.
- Despliegue en la nube (Render o Railway).

---

*TDI-601 Aplicaciones Web II · ULEAM · 2026-1*
