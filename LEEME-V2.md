# Teldrive v2.0.0-drive-es.10

Integra las actualizaciones del original hasta `67bf180` y reúne las mejoras de bots y Desktop en una sola versión.

- Bots identificados y ordenados por **nombre, @alias e ID**, con creación desde un gestor y recuperación de operaciones interrumpidas.
- WebView2 usa el host de ventana nativo para evitar el error gráfico `0x8876086C` al iniciar o restaurar desde bandeja.
- Menús de bandeja: las acciones se ejecutan al terminar el menú nativo, evitando la reentrada en WPF/WebView2; texto del icono limitado al tamaño de Windows.
- **Pantalla completa** en vídeo web y reproductor nativo. En el nativo: botón, F/F11, doble clic y Esc, conservando la ventana anterior.
- **Desarrollador · estadísticas** en audio/vídeo: datos de reproducción, pistas/códecs y búfer. Desktop nativo añade resolución, FPS, aceleración, tasas de bits, sincronización A/V, fotogramas perdidos y otros datos reales de mpv por IPC.
- Telegram: cuenta real que recibe cada fragmento, nombre/@alias/ID, velocidad móvil de 2 s, bytes, solicitudes activas/total, errores, latencia y último canal/mensaje/offset. Las cifras suman las lecturas del archivo, incluidas búsquedas y otras reproducciones simultáneas; una velocidad cero puede indicar reproducción desde el búfer. Datos visibles para el propietario, sin claves/tokens/URL de reproducción.
- Del original: provisión independiente de bots, gestión de carga y cuentas revocadas, favoritos/recientes, preferencias por cuenta, tamaños de carpetas, orden natural, búsqueda/lectores y descargas compartidas directas con autorización por cookie. Se conservan Spam, temas, diez idiomas, tareas, rclone y transmisión.

## Actualización

Cierra Desktop, conserva tu carpeta de datos y configuración, sustituye el ejecutable y vuelve a abrirlo. Las migraciones propias 8–10 se conservan; las nuevas pasan a 11–16. No regeneres claves ni repitas la migración. El paquete Windows incluye .NET, PostgreSQL, rclone, FFmpeg/FFprobe y mpv. Los CLI para Windows/Linux/macOS x64 y ARM64 necesitan PostgreSQL y FFmpeg externos. Android/iOS conservan sus clientes y usan la UI del servidor actualizado.

## Comprobaciones y límites

Go unitario completo, vet, TypeSpec, TypeScript, lint y compilación de UI; regresiones de bots, búsqueda, preferencias, favoritos, descargas y estadísticas con API simulada. Prueba nativa de reproducción real con mpv, IPC, pausa, pantalla completa/restauración, panel de desarrollador y callbacks de bandeja. Arranque de PostgreSQL/WebView2 y actualización de una base anterior en carpetas aisladas. Los manifiestos incluyen las pruebas concretas ejecutadas y los SHA256.

No se han usado bots ni credenciales reales de Telegram. El harness de integración PostgreSQL no pudo ejecutarse porque Docker no está disponible; las pruebas de autorización de estadísticas y reconciliación del esquema original se compilaron, pero necesitan ese harness para ejecutarse. Linux/macOS/ARM64 se compilan y verifican por arquitectura/metadatos, sin ejecución nativa. La prueba iPhone usa una ruta explícita al WebKit instalado. Los datos de mpv solo se muestran cuando el motor los proporciona; otros orígenes rclone mantienen los datos del reproductor sin atribuir cuentas de Telegram.
