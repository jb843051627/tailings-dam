package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

func TestTD10_NextInspectionDateWrong(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td10_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	ctx := context.Background()
	dam, err := st.CreateDam(ctx, &model.Dam{
		Name:        "Test Dam",
		Location:    "Test Location",
		Latitude:     30.0,
		Longitude:    120.0,
		Capacity:     1000000,
		HazardLevel: model.HazardLevelMedium,
		Status:      model.DamStatusActive,
	})
	if err != nil {
		t.Fatalf("failed to create dam: %v", err)
	}

	// Create inspection with zero scheduled date
	_, err = st.CreateInspection(ctx, &model.Inspection{
		DamID:         dam.ID,
		Inspector:     "Test Inspector",
		Title:         "Zero Date Inspection",
		ScheduledDate: time.Time{}, // zero value
		Status:        model.InspectionStatusPending,
		Priority:      model.InspectionPriorityNormal,
	})
	if err != nil {
		t.Fatalf("failed to create inspection: %v", err)
	}

	svc := NewInspectionService(st)

	// List inspections - zero-time scheduled date should not be marked as overdue
	inspections, err := svc.ListInspections(ctx)
	if err != nil {
		t.Fatalf("failed to list inspections: %v", err)
	}

	for _, insp := range inspections {
		if insp.ScheduledDate.IsZero() && insp.Status == model.InspectionStatusOverdue {
			t.Fatal("inspection with zero scheduled date was marked as overdue, expected to remain pending")
		}
	}
}
