# Teldrive v2.0.0-drive-es.13

Desktop can now connect to an existing central Go or Node server using its TOML:

```toml
[desktop]
server-url = "https://cloud.webrom.lat"
```

Existing local TOML files keep working. Remote mode does not start PostgreSQL,
the Go server, migrations or workers on the PC. Closing Desktop disconnects the
client and leaves the central server running. The client needs no Supabase
credentials or installation encryption keys. Login and account permissions are
still enforced by the selected server.

Install the Windows package in a new folder. For portable remote use, create
`Teldrive.portable` next to the executable and save the TOML above as
`TeldriveData/config.toml`. Existing installations can use the advanced TOML
editor; back up their complete local configuration before replacing it. The
separate `desktop/REMOTE.md` in the sources documents switching modes and API
compatibility. Keep a server-generated API key for each rclone connection;
Go and Node sessions/API keys are separate.

Remote servers must use HTTPS, except for loopback development. URLs containing
credentials, non-root paths, queries or fragments are rejected. Readiness checks
do not follow redirects; browser navigation and native messages remain pinned to
the configured origin. Native integrations retain the existing owner-role check.

The web interface and supported operations come from the selected server. Node
shares the v2 catalog but does not implement every Go capability, including
River jobs. This release does not migrate or replace either server. Performance
depends on the server, network path and media conversion; no comparative transfer
speed claim is made.

Windows Desktop version: 2.0.0.13. The compatible rclone and mobile applications
are unchanged. See VALIDACION.json for executed checks and platform limitations.

Exact source revisions: `source-revisions.json`.
