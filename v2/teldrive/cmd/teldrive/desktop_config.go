package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
)

// A remote Desktop configuration belongs to the client, not the local server.
// Parse it before loading server settings so clients need no database secrets.
func desktopServerURL(path string) (string, error) {
	if path == "" || !strings.EqualFold(filepath.Ext(path), ".toml") {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read desktop configuration: %w", err)
	}
	settings, err := toml.Parser().Unmarshal(data)
	if err != nil {
		// Parser errors may include private lines from an existing server TOML.
		return "", fmt.Errorf("invalid TOML in desktop configuration")
	}
	value, exists := settings["desktop"]
	if !exists {
		return "", nil
	}
	desktop, ok := value.(map[string]any)
	if !ok {
		return "", fmt.Errorf("desktop must be a TOML table")
	}
	value, exists = desktop["server-url"]
	if !exists {
		return "", nil
	}
	server, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("desktop.server-url must be a string")
	}
	if server == "" {
		return "", nil
	}
	u, err := url.Parse(server)
	if err != nil || u.Hostname() == "" || u.Opaque != "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.RawFragment != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || strings.Contains(server, "#") {
		return "", fmt.Errorf("desktop.server-url must be an HTTPS origin without credentials, path, query or fragment")
	}
	loopback := strings.EqualFold(u.Hostname(), "localhost")
	if ip := net.ParseIP(u.Hostname()); ip != nil {
		loopback = ip.IsLoopback()
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && loopback) {
		return "", fmt.Errorf("desktop.server-url requires HTTPS (HTTP is permitted only on loopback)")
	}
	// url.Parse accepts an empty explicit port and ports beyond the TCP range.
	if strings.HasSuffix(u.Host, ":") || (u.Port() != "" && !validDesktopPort(u.Port())) {
		return "", fmt.Errorf("invalid desktop.server-url port")
	}
	u.Path = "/"
	return u.String(), nil
}

func validDesktopPort(port string) bool {
	n := 0
	for _, digit := range port {
		if digit < '0' || digit > '9' {
			return false
		}
		n = n*10 + int(digit-'0')
		if n > 65535 {
			return false
		}
	}
	return n > 0
}
