package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDesktopRemoteInfoNeedsNoDatabaseOrCredentials(t *testing.T) {
	for _, server := range []string{"https://cloud.example.com", "https://cloud.example.com/", "https://cloud.example.com:8443", "http://127.0.0.1:3000", "http://[::1]:3000", "http://localhost:3000"} {
		t.Run(server, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.toml")
			// Existing server fields may remain in a shared TOML. Remote mode
			// must not initialize/validate that server or print its secrets.
			data := "[desktop]\nserver-url = \"" + server + "\"\n[database]\nurl = \"private-unused-invalid-url\"\n"
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			cmd := newRootCommand()
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetArgs([]string{"check", "--config", path, "--desktop-info"})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			var info map[string]string
			if err := json.Unmarshal(out.Bytes(), &info); err != nil {
				t.Fatal(err)
			}
			if len(info) != 1 || info["serverUrl"] != strings.TrimSuffix(server, "/")+"/" || strings.Contains(out.String(), "private-unused") {
				t.Fatalf("unexpected remote info: %s", out.String())
			}
		})
	}
}

func TestDesktopRemoteConfigRejectsUnsafeOrigins(t *testing.T) {
	for _, server := range []string{"http://cloud.example.com", "ftp://cloud.example.com", "https://user:private-secret@cloud.example.com", "https://cloud.example.com/path", "https://cloud.example.com/?private-secret", "https://cloud.example.com/#private-secret", "https://cloud.example.com/?", "https://cloud.example.com/#", "https://cloud.example.com/%2f", "https://cloud.example.com:0", "https://cloud.example.com:65536", "https://cloud.example.com:", "//cloud.example.com"} {
		t.Run(server, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(path, []byte("[desktop]\nserver-url = \""+server+"\"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := desktopServerURL(path)
			if err == nil || strings.Contains(err.Error(), "private-secret") {
				t.Fatalf("expected redacted error, got %v", err)
			}
		})
	}
}

func TestDesktopRemoteConfigPreservesLocalAndRejectsInvalidTOML(t *testing.T) {
	for _, data := range []string{"[desktop]\nserver-url = 123", "[desktop]\nserver-url = [\"https://example.com\"]", "desktop = 123", "[desktop]\nserver-url = \"private-secret"} {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := desktopServerURL(path)
		if err == nil || strings.Contains(err.Error(), "private-secret") {
			t.Fatalf("expected redacted configuration error, got %v", err)
		}
	}
	for _, data := range []string{"[http]\naddress = \"127.0.0.1:8080\"", "[desktop]\nserver-url = \"\""} {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		server, err := desktopServerURL(path)
		if err != nil || server != "" {
			t.Fatalf("local configuration changed: %q, %v", server, err)
		}
	}
}
