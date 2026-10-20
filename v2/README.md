# Teldrive v2.0.0-drive-es.12 — conexión con bases ya migradas

Corrige el error `legacy TelDrive schema detected: database.auto-migrate-legacy is disabled` cuando la base de datos ya usa el esquema v2, pero reaparece `public.goose_db_version`.

El servidor comprueba la estructura de las tablas antes de clasificar la base. Una base identificada como v2 continúa con su arranque normal y no entra otra vez en el migrador de v1, incluso si la migración automática está habilitada. Las bases realmente antiguas, incompletas o mixtas conservan las protecciones existentes.

No necesitas activar `database.auto-migrate-legacy`, volver a importar archivos ni generar claves nuevas para corregir este caso. La versión no añade migraciones de datos ni modifica tu configuración. Conserva las mejoras de bots, reproducción, estadísticas y la corrección de reapertura de .11.

## Actualización

Desde el icono de la bandeja, elige **Salir y detener servicios**, descomprime el paquete nuevo y sustituye el ejecutable. Conserva tu carpeta de datos, configuración y claves. En modo portátil, conserva `TeldriveData` junto al ejecutable. Los CLI se recompilan para Windows, Linux y macOS, x64 y ARM64. Rclone y los clientes móviles conservan sus versiones anteriores.

El registro de Actividad mantiene su fondo oscuro y texto legible independientemente del tema o foco de Windows.

## Comprobaciones

La regresión se reprodujo con el código anterior en PostgreSQL aislado. Pasaron las suites de integración de base de datos y migración v1 mediante el harness del proyecto con Docker/PostgreSQL 18: marcador antiguo sobre v2, migración automática habilitada/deshabilitada, rechazo del migrador explícito sobre v2 y protección de esquemas antiguos, incompletos o mixtos. Se verificó que se conservan el esquema, los registros y el marcador original.

También se comprueban Go unitario completo, vet, compilación de la UI, metadatos de los seis CLI, equivalencia del servidor integrado en Desktop, arranque nativo de PostgreSQL/WebView2, actualización conservando configuración y las siete pruebas de reapertura. Los resultados y SHA256 concretos se publican en los manifiestos de la release.

La UI funcional, las API y los clientes móviles no cambian. Los resultados de las 57 regresiones de navegador de .10 siguen siendo el antecedente de esa UI; no se presentan como una nueva ejecución en .12. No se realizan transferencias reales con Telegram ni se afirma ejecución nativa de Linux/macOS/ARM64 en esta publicación.

Fuentes: `teldrive/` y `rclone/`; revisiones exactas en `source-revisions.json`.
