//go:build integration

package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	api "github.com/tgdrive/teldrive/v2/internal/api"
	"github.com/tgdrive/teldrive/v2/internal/authn"
	"github.com/tgdrive/teldrive/v2/internal/catalog"
	"github.com/tgdrive/teldrive/v2/internal/health"
	"github.com/tgdrive/teldrive/v2/internal/secureblob"
	testpostgres "github.com/tgdrive/teldrive/v2/internal/testutil/postgres"
	"github.com/tgdrive/teldrive/v2/internal/uploads"
)

func TestProfileReportsPublicURLOnlyWhenConfigured(t *testing.T) {
	db := testpostgres.New(t)
	ctx := context.Background()
	if _, err := db.Pool.Exec(ctx, "INSERT INTO users (user_id) VALUES (1001)"); err != nil {
		t.Fatal(err)
	}
	cipher, err := secureblob.New("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	if err != nil {
		t.Fatal(err)
	}
	authService, err := authn.NewService(db.Pool, cipher, unusedTelegramLogin{}, authn.Config{
		SigningKey: "0123456789abcdef0123456789abcdef", Issuer: "test",
		AccessTokenTTL: time.Hour, RefreshTokenTTL: time.Hour, LoginFlowTTL: time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	handler := api.NewHandler(catalog.NewService(db.Pool, nil), uploads.NewService(db.Pool), nil, nil, health.NewService("test", db.Pool), 0, nil).
		ConfigureDomains(authService, nil, nil, nil, nil, nil)
	server, err := api.NewServer(handler, api.NewSecurity(tokenUserAuthenticator{"owner-token": 1001}))
	if err != nil {
		t.Fatal(err)
	}
	profile := func() map[string]any {
		t.Helper()
		response := performRequest(t, server, http.MethodGet, "/v1/me", nil, map[string]string{"Authorization": "Bearer owner-token"})
		if response.Code != http.StatusOK {
			t.Fatalf("GET /v1/me status = %d body = %s", response.Code, response.Body.String())
		}
		var decoded map[string]any
		if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
			t.Fatalf("decode profile: %v", err)
		}
		return decoded
	}

	if _, found := profile()["publicUrl"]; found {
		t.Fatalf("publicUrl reported without a configured public address")
	}
	handler.PublicURL = url.URL{Scheme: "https", Host: "tel.example"}
	if got := profile()["publicUrl"]; got != "https://tel.example" {
		t.Fatalf("publicUrl = %v, want https://tel.example", got)
	}
}
