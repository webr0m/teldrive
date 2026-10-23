# Connect Desktop to an existing Teldrive server

Desktop can keep running its bundled local server, or connect to a remote Go or
Node server. Existing TOML files and local databases keep their current behavior.

To use the central server, save this in Desktop's `config.toml`:

```toml
[desktop]
server-url = "https://cloud.webrom.lat"
```

The file lives in `%LOCALAPPDATA%\TeldriveV2\config.toml`, or `TeldriveData\config.toml`
next to the executable in portable mode. A fresh portable installation needs only
this TOML and the `Teldrive.portable` marker; Desktop creates its client state.
Existing installations can use the advanced TOML editor, save, and reconnect.

The URL must be an HTTPS origin, without credentials, paths, query parameters or
fragments. HTTP is accepted only for a loopback development server. Desktop checks
the server's readiness without following redirects and keeps browser navigation
and native messages restricted to the configured origin.

Remote mode does not start PostgreSQL, the Go server, migrations or background
workers on the PC. Disconnecting or closing Desktop leaves the central server
running. Log in with your Telegram account on that server; its allowlist and roles
remain authoritative. The TOML does not grant access to another user's files.

The remote client does not need Supabase credentials, Telegram application secrets
or database/content encryption keys. Those stay on the central server. Browser
sessions, rclone configuration, cache and Windows integration settings remain local.
Local rclone connections must use an API key created on the selected server;
existing API keys and browser sessions are not interchangeable between Go and Node.

Playback, catalog operations and features follow the remote server's API and web
interface. Node supports the shared v2 catalog, but does not implement every Go
feature (for example, Go's River jobs). Remote mode does not migrate that data or
pretend that the two backends implement identical capabilities. Native integrations
retain Desktop's existing owner-role check.

To return to a previous local installation, keep a backup of its complete TOML.
Remove `desktop.server-url` (or set it to an empty string) and restore the local
server configuration. Desktop preserves its previous local database settings when
switching modes; it never deletes the local database or regenerates existing keys.
