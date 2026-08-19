package store

import (
	"context"
	"testing"
	"time"

	"tailings-dam/internal/model"
)

func TestTD07_BatchCreateInspectionsRollbackOnError(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	dam := mustCreateDam(t, st)

	_, err := st.db.Exec("CREATE UNIQUE INDEX idx_insp_unique ON inspections(dam_id, title)")
	if err != nil {
		t.Fatalf("failed to create unique index: %v", err)
	}

	now := time.Now()
	inspections := []*model.Inspection{
		{
			DamID:         dam.ID,
			Inspector:     "Inspector A",
			Title:         "Duplicate Title",
			ScheduledDate: now,
			Status:        model.InspectionStatusPending,
			Priority:      model.InspectionPriorityNormal,
		},
		{
			DamID:         dam.ID,
			Inspector:     "Inspector B",
			Title:         "Duplicate Title",
			ScheduledDate: now,
			Status:        model.InspectionStatusPending,
			Priority:      model.InspectionPriorityNormal,
		},
		{
			DamID:         dam.ID,
			Inspector:     "Inspector C",
			Title:         "Unique Title",
			ScheduledDate: now,
			Status:        model.InspectionStatusPending,
			Priority:      model.InspectionPriorityNormal,
		},
	}

	_, err = st.BatchCreateInspections(ctx, inspections)
	if err == nil {
		t.Fatal("expected error for batch with duplicate title, got nil")
	}

	count, err := st.GetInspectionCount(ctx)
	if err != nil {
		t.Fatalf("failed to count inspections: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 inspections after rollback, got %d", count)
	}
}

func TestTD07_BatchCreateAlertsRollbackOnError(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	dam := mustCreateDam(t, st)

	_, err := st.db.Exec("CREATE UNIQUE INDEX idx_alert_unique ON alerts(dam_id, title)")
	if err != nil {
		t.Fatalf("failed to create unique index: %v", err)
	}

	alerts := []*model.Alert{
		{
			DamID:   dam.ID,
			Level:   model.AlertLevelWarning,
			Status:  model.AlertStatusActive,
			Title:   "Duplicate Alert",
			Message: "first",
		},
		{
			DamID:   dam.ID,
			Level:   model.AlertLevelWarning,
			Status:  model.AlertStatusActive,
			Title:   "Duplicate Alert",
			Message: "second",
		},
		{
			DamID:   dam.ID,
			Level:   model.AlertLevelWarning,
			Status:  model.AlertStatusActive,
			Title:   "Unique Alert",
			Message: "third",
		},
	}

	_, err = st.BatchCreateAlerts(ctx, alerts)
	if err == nil {
		t.Fatal("expected error for batch with duplicate title, got nil")
	}

	count, err := st.GetActiveAlertCount(ctx)
	if err != nil {
		t.Fatalf("failed to count alerts: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 alerts after rollback, got %d", count)
	}
}
