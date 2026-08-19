package store

import (
	"context"
	"sync"
	"testing"

	"tailings-dam/internal/model"
)

func TestTD03_ConcurrentUpdateDataRace(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	dam := mustCreateDam(t, st)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, err := st.CreateMonitoringPoint(ctx, &model.MonitoringPoint{
				DamID: dam.ID,
				Name:  "Concurrent Point",
				Code:  "CP",
				Type:  model.PointTypeSeepage,
				Status: model.PointStatusActive,
			})
			if err != nil {
				t.Errorf("concurrent create failed: %v", err)
			}
		}(i)
	}
	wg.Wait()
}
