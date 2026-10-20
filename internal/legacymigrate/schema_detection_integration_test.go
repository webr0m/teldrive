//go:build integration

package legacymigrate_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/tgdrive/teldrive/v2/internal/database"
	"github.com/tgdrive/teldrive/v2/internal/legacymigrate"
	testpostgres "github.com/tgdrive/teldrive/v2/internal/testutil/postgres"
)

func TestMigratedDatabaseWithRecreatedGooseMarker(t *testing.T) {
	db := testpostgres.New(t)
	ctx := context.Background()
	if _, err := db.Pool.Exec(ctx, `
CREATE TABLE public.goose_db_version(id bigint PRIMARY KEY);
INSERT INTO public.goose_db_version VALUES (73);
INSERT INTO teldrive.users(user_id, display_name) VALUES (101, 'Preserved owner');
INSERT INTO teldrive.files(id, user_id, name, kind, mod_time)
VALUES ('00000000-0000-4000-8000-000000000001', 101, 'Preserved folder', 'folder', now());
`); err != nil {
		t.Fatal(err)
	}
	var schemaOID uint32
	if err := db.Pool.QueryRow(ctx, `SELECT oid FROM pg_namespace WHERE nspname='teldrive'`).Scan(&schemaOID); err != nil {
		t.Fatal(err)
	}
	for _, automatic := range []bool{false, true} {
		t.Run(fmt.Sprintf("automatic=%v", automatic), func(t *testing.T) {
			cfg := database.Config{URL: db.URL, AutoMigrateLegacy: automatic}
			// No data key or verifier is needed: a confirmed v2 database must
			// never enter the legacy copier, even when automatic migration is on.
			if _, migrated, err := legacymigrate.MigrateIfNeeded(ctx, cfg, "", nil); err != nil || migrated {
				t.Fatalf("already migrated database: migrated=%v, error=%v", migrated, err)
			}
			if err := database.Migrate(ctx, cfg); err != nil {
				t.Fatalf("normal v2 startup migrations: %v", err)
			}
		})
	}
	_, err := legacymigrate.Run(ctx, legacymigrate.Config{
		SourceURL: db.URL, Target: database.Config{URL: db.URL},
		DataKey: strings.Repeat("A", 43), EncryptionKeyVersion: 1, Apply: true,
	})
	if err == nil || !strings.Contains(err.Error(), "not a legacy") {
		t.Fatalf("explicit legacy copier must reject v2: %v", err)
	}
	var finalOID uint32
	var owner, folder string
	var marker, extraSchemas int
	if err := db.Pool.QueryRow(ctx, `SELECT
 (SELECT oid FROM pg_namespace WHERE nspname='teldrive'),
 (SELECT display_name FROM teldrive.users WHERE user_id=101),
 (SELECT name FROM teldrive.files WHERE user_id=101),
 (SELECT count(*) FROM public.goose_db_version WHERE id=73),
 (SELECT count(*) FROM pg_namespace WHERE nspname LIKE 'teldrive_legacy_backup_%' OR nspname LIKE 'teldrive_v2_staging_%')
`).Scan(&finalOID, &owner, &folder, &marker, &extraSchemas); err != nil {
		t.Fatal(err)
	}
	if finalOID != schemaOID || owner != "Preserved owner" || folder != "Preserved folder" || marker != 1 || extraSchemas != 0 {
		t.Fatalf("existing schemas or data changed: oid=%v owner=%q folder=%q marker=%d extra schemas=%d", finalOID, owner, folder, marker, extraSchemas)
	}
}

func TestLegacyGuardStillRejectsUnconfirmedV2Schemas(t *testing.T) {
	for _, shape := range []string{"marker only", "legacy files", "mixed files"} {
		t.Run(shape, func(t *testing.T) {
			db := testpostgres.New(t)
			ctx := context.Background()
			if shape == "marker only" {
				if _, err := db.Pool.Exec(ctx, `DROP SCHEMA teldrive CASCADE`); err != nil {
					t.Fatal(err)
				}
			} else if shape == "legacy files" {
				if _, err := db.Pool.Exec(ctx, `DROP SCHEMA teldrive CASCADE; CREATE SCHEMA teldrive; CREATE TABLE teldrive.files(id uuid, type text, parts jsonb)`); err != nil {
					t.Fatal(err)
				}
			} else {
				if _, err := db.Pool.Exec(ctx, `ALTER TABLE teldrive.files ADD COLUMN type text, ADD COLUMN parts jsonb`); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := db.Pool.Exec(ctx, `CREATE TABLE public.goose_db_version(id bigint)`); err != nil {
				t.Fatal(err)
			}
			if err := database.Migrate(ctx, database.Config{URL: db.URL}); !errors.Is(err, database.ErrLegacySchema) {
				t.Fatalf("legacy schema was not blocked: %v", err)
			}
			if _, migrated, err := legacymigrate.MigrateIfNeeded(ctx, database.Config{URL: db.URL}, "", nil); migrated || err == nil || !strings.Contains(err.Error(), "security.data-key") {
				t.Fatalf("legacy detection was bypassed: migrated=%v, error=%v", migrated, err)
			}
		})
	}
}
