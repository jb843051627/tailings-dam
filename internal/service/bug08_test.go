package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

func TestTD08_ReportSortPollutesStore(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "td08_test.db")
	st, err := store.New(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer st.Close()
	defer os.Remove(dbPath)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_, err := st.CreateDam(ctx, &model.Dam{
			Name:        "Dam " + string(rune('A'+i)),
			Location:    "Location",
			Latitude:     30.0,
			Longitude:    120.0,
			Capacity:     1000000,
			HazardLevel: model.HazardLevelMedium,
			Status:      model.DamStatusActive,
		})
		if err != nil {
			t.Fatalf("failed to create dam: %v", err)
		}
	}

	svc := NewDamService(st)

	first, err := svc.ListDams(ctx)
	if err != nil {
		t.Fatalf("first ListDams failed: %v", err)
	}
	if len(first) != 5 {
		t.Fatalf("expected 5 dams, got %d", len(first))
	}

	// Reverse the first result
	for i, j := 0, len(first)-1; i < j; i, j = i+1, j-1 {
		first[i], first[j] = first[j], first[i]
	}

	// Second call - should return original order if cache not polluted
	second, err := svc.ListDams(ctx)
	if err != nil {
		t.Fatalf("second ListDams failed: %v", err)
	}

	if first[0].ID == second[0].ID {
		t.Fatal("store cache polluted by sort: second call returns reversed order")
	}
}
