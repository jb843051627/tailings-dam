package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"tailings-dam/internal/store"
)

func newTestStoreSvc(t *testing.T) *store.Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "svc_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	t.Cleanup(func() {
		st.Close()
		os.Remove(dbPath)
	})
	return st
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
