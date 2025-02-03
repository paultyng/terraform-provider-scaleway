package rdb_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/terraform-provider-scaleway/v2/internal/services/rdb"
)

func TestPrivilegeV1SchemaUpgradeFunc(t *testing.T) {
	v0Schema := map[string]interface{}{
		"id":            "fr-par/11111111-1111-1111-1111-111111111111",
		"region":        "fr-par",
		"database_name": "database",
		"user_name":     "username",
	}
	v1Schema := map[string]interface{}{
		"id":            "fr-par/11111111-1111-1111-1111-111111111111/database/username",
		"region":        "fr-par",
		"database_name": "database",
		"user_name":     "username",
	}

	actual, err := rdb.PrivilegeV1SchemaUpgradeFunc(context.Background(), v0Schema, nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(v1Schema, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", v1Schema, actual)
	}
}

func TestExtractEngineVersion(t *testing.T) {
	tests := []struct {
		engine         string
		expected       int
		expectingError bool
	}{
		{"postgresql-15", 15, false},
		{"mysql-8.0", 8, false},
		{"redis-6", 6, false},
		{"mariadb-10.5", 10, false}, // Only extracts the major version
		{"invalid-engine", 0, true}, // No version to extract
		{"", 0, true},               // Empty string case
		{"mongodb-3.6", 3, false},   // Extracts only the major version
	}

	for _, tt := range tests {
		t.Run(tt.engine, func(t *testing.T) {
			result, err := rdb.ExtractEngineVersion(tt.engine)

			if tt.expectingError {
				if err == nil {
					t.Errorf("expected an error for engine %q, but got none", tt.engine)
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error for engine %q, but got: %s", tt.engine, err)
				}

				if result != tt.expected {
					t.Errorf("expected version %d for engine %q, but got %d", tt.expected, tt.engine, result)
				}
			}
		})
	}
}
