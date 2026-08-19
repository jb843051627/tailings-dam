package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/store"
)

func TestTD06_DrainageNilPanic(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td06_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	ctx := context.Background()

	// Non-existent drainage system should return error, not (nil, nil)
	d, err := st.GetDrainageSystem(ctx, 99999)
	if d == nil && err == nil {
		t.Fatal("GetDrainageSystem returned (nil, nil) for non-existent record, expected error")
	}
}
