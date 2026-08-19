package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/store"
)

func TestTD01_RecordReadingNilPointPanics(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td01_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	ctx := context.Background()
	_ = ctx

	// Non-existent monitoring point should return error, not (nil, nil)
	mp, err := st.GetMonitoringPoint(ctx, 99999)
	if mp == nil && err == nil {
		t.Fatal("GetMonitoringPoint returned (nil, nil) for non-existent point, expected error")
	}
}
