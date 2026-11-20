# Teldrive v2.0.0-drive-es.14

Links that the web interface copies for other people now use the public origin
configured on the server instead of the address in the browser bar. Servers
behind a reverse proxy keep reporting the reachable domain:

```toml
[http]
address = "127.0.0.1:8086"
public-url = "https://tel.example.com"
```

`GET /api/v1/me` returns that value as `publicUrl` only when it is configured,
so installations that leave it empty keep building links from their own origin
and nothing changes for them. Right-click -> Herramientas -> "Copiar enlace de
descarga" and the selection toolbar now resolve that origin, which fixes links
that only reached a private LAN address such as `http://192.168.1.50:8086`.
"Crear enlace" share links were already absolute and behave as before.

The settings information screen shows which public address the server is
advertising, so operators can confirm the configuration without reading logs.

Download links still require a logged-in session at that address. For access
without an account, keep using "Crear enlace público", which returns a token
link handled by the share endpoints.

The Desktop remote documentation and the remote TOML example no longer point at
the retired hosting domain; they now use a neutral placeholder server address.
Remote mode requirements are unchanged: HTTPS origin, no credentials, path,
query or fragment.

Windows Desktop version: 2.0.0.14. The compatible rclone and mobile applications
are unchanged. See VALIDACION.json for executed checks and platform limitations.

Exact source revisions: `source-revisions.json`.
