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

func TestTD05_BatchIngestIgnoresContext(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td05_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	svc := NewReadingService(st)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	batch := &model.BatchReadingInput{
		SeepageReadings: []model.SeepageReadingInput{
			{PointID: 1, FlowRate: 10.0, Turbidity: 5.0, PH: 7.0, Temperature: 20.0},
		},
	}

	_, err = svc.BatchIngest(ctx, batch)
	if err == nil {
		t.Fatal("BatchIngest should return error when context is cancelled, got nil")
	}

	// The error should preserve the context.Canceled chain
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected error to wrap context.Canceled, got: %v", err)
	}
}
