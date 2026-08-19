package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	t.Cleanup(func() { st.Close(); os.Remove(dbPath) })
	return st
}

func mustCreateDam(t *testing.T, s *Store) *model.Dam {
	t.Helper()
	ctx := context.Background()
	dam, err := s.CreateDam(ctx, &model.Dam{
		Name:        "Test Dam",
		Location:    "Test Location",
		Province:    "Test Province",
		Latitude:     30.0,
		Longitude:    120.0,
		Capacity:     1000000,
		HazardLevel: model.HazardLevelMedium,
		Status:      model.DamStatusActive,
	})
	if err != nil {
		t.Fatalf("failed to create dam: %v", err)
	}
	return dam
}

func mustCreateMonitoringPoint(t *testing.T, s *Store, damID int64) *model.MonitoringPoint {
	t.Helper()
	ctx := context.Background()
	mp, err := s.CreateMonitoringPoint(ctx, &model.MonitoringPoint{
		DamID: damID,
		Name:  "Test Point",
		Code:  "TP001",
		Type:  model.PointTypeSeepage,
		Status: model.PointStatusActive,
	})
	if err != nil {
		t.Fatalf("failed to create monitoring point: %v", err)
	}
	return mp
}

var _ = context.Background
