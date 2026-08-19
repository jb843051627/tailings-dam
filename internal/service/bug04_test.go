package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

func TestTD04_ErrorWrappingBreaksErrorsIs(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td04_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	svc := NewAlertService(st)
	ctx := context.Background()

	// Create alert with invalid input (empty dam_id)
	_, err = svc.CreateAlert(ctx, &model.AlertInput{
		DamID:    0,
		Level:    model.AlertLevelWarning,
		Message:  "test",
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	// The validation error should be wrapped; check errors.Is
	var valErr *model.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected errors.As to find ValidationError, got: %v", err)
	}
}
