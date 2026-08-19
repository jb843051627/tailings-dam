package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

func TestTD09_ActivateIgnoresValidation(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td09_test.db")
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

	insp, err := st.CreateInspection(ctx, &model.Inspection{
		DamID:         dam.ID,
		Inspector:     "Test Inspector",
		Title:         "Test Inspection",
		ScheduledDate: mustParseTime("2026-06-01T00:00:00Z"),
		Status:        model.InspectionStatusPending,
		Priority:      model.InspectionPriorityNormal,
	})
	if err != nil {
		t.Fatalf("failed to create inspection: %v", err)
	}

	svc := NewInspectionService(st)

	// Invalid input: empty inspector and title
	invalidInput := &model.InspectionInput{
		DamID:         dam.ID,
		Inspector:     "",
		Title:         "",
		ScheduledDate: "2026-06-01T00:00:00Z",
		Priority:      model.InspectionPriorityNormal,
	}

	_, err = svc.ActivateInspection(ctx, insp.ID, invalidInput)
	if err == nil {
		t.Fatal("ActivateInspection should return validation error for invalid input, got nil")
	}
}
