package store

import (
	"context"
	"testing"

	"tailings-dam/internal/model"
)

func TestTD02_ListAlertsByDamCachePollution(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	dam := mustCreateDam(t, st)

	for i := 0; i < 3; i++ {
		_, err := st.CreateAlert(ctx, &model.Alert{
			DamID:   dam.ID,
			Level:   model.AlertLevelWarning,
			Status:  model.AlertStatusActive,
			Title:   "Original Title",
			Message: "test alert",
		})
		if err != nil {
			t.Fatalf("failed to create alert: %v", err)
		}
	}

	first, err := st.ListAlertsByDam(ctx, dam.ID)
	if err != nil {
		t.Fatalf("first ListAlertsByDam failed: %v", err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(first))
	}

	first[0].Title = "MUTATED"

	second, err := st.ListAlertsByDam(ctx, dam.ID)
	if err != nil {
		t.Fatalf("second ListAlertsByDam failed: %v", err)
	}

	if second[0].Title == "MUTATED" {
		t.Fatal("cache polluted: second call returns mutated data from first call's slice")
	}
}

func TestTD02_ListAlertsByStatusCachePollution(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	dam := mustCreateDam(t, st)

	for i := 0; i < 3; i++ {
		_, err := st.CreateAlert(ctx, &model.Alert{
			DamID:   dam.ID,
			Level:   model.AlertLevelWarning,
			Status:  model.AlertStatusActive,
			Title:   "Original Title",
			Message: "test alert",
		})
		if err != nil {
			t.Fatalf("failed to create alert: %v", err)
		}
	}

	first, err := st.ListAlertsByStatus(ctx, model.AlertStatusActive)
	if err != nil {
		t.Fatalf("first ListAlertsByStatus failed: %v", err)
	}
	if len(first) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(first))
	}

	first[0].Message = "POLLUTED"

	second, err := st.ListAlertsByStatus(ctx, model.AlertStatusActive)
	if err != nil {
		t.Fatalf("second ListAlertsByStatus failed: %v", err)
	}

	if second[0].Message == "POLLUTED" {
		t.Fatal("cache polluted: second call returns mutated data from first call's slice")
	}
}
