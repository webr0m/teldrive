# Teldrive Desktop 2.0.0.11 — corrección de reapertura

Corrige el error «No se puede mostrar Window cuando ShowActivated es false y WindowState está establecido en Maximized» al volver a abrir Desktop después de cerrarlo en la bandeja.

La ventana se activa antes de mostrarse y conserva el estado maximizado. Se cubren la restauración desde la bandeja y abrir el ejecutable otra vez cuando ya sigue activo. Minimizar y restaurar devuelve la ventana a su estado normal.

## Actualización

En el icono de Teldrive de la bandeja, elige **Salir y detener servicios**. Descomprime el paquete nuevo y sustituye `Teldrive Desktop v2.exe`. Conserva tu carpeta de datos, configuración y claves. No ejecutes la configuración inicial otra vez. Si utilizas modo portátil, mantén `TeldriveData` junto al ejecutable.

Este es un hotfix de Windows Desktop: el servidor, UI y componentes integrados conservan la versión .10. Los seis paquetes CLI y los clientes móviles publicados anteriormente siguen siendo válidos y no se recompilan para este cambio. Las mejoras de bots, pantalla completa y estadísticas de .10 se conservan.

## Comprobaciones

Se reprodujo el error original con la implementación anterior. El ejecutable corregido se prueba en estado normal y maximizado, al cerrar y restaurar directamente, desde el menú de bandeja y lanzando un segundo proceso real, además de minimizar y restaurar. Se comprueban también el arranque de PostgreSQL/WebView2, pantalla completa y la reapertura desde un inicio en bandeja con una base de pruebas existente, conservando su configuración. Los resultados concretos se incluyen en `VALIDACION.json`.

No se usan credenciales reales de Telegram. Este hotfix no añade migraciones ni modifica el servidor. El manifiesto identifica por separado la revisión de Desktop y la revisión del servidor integrado; los archivos y sus tamaños/SHA256 se verifican después de subirlos.

Fuentes: `teldrive/` y `rclone/`; revisiones exactas en `source-revisions.json`.
