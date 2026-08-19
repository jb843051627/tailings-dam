package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateSeepageReading 创建渗流读数
func (s *Store) CreateSeepageReading(ctx context.Context, r *model.SeepageReading) (*model.SeepageReading, error) {
	r.CreatedAt = time.Now()
	if r.RecordedAt.IsZero() {
		r.RecordedAt = time.Now()
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO seepage_readings (point_id, flow_rate, turbidity, ph, temperature,
			recorded_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		r.PointID, r.FlowRate, r.Turbidity, r.PH, r.Temperature,
		r.RecordedAt, r.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create seepage reading: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get seepage reading id: %v", err)
	}
	r.ID = id

	// 更新监测点的最后读数时间
	_ = s.UpdateLastReading(ctx, r.PointID, r.RecordedAt)

	return r, nil
}

// ListSeepageReadings 列出渗流读数
func (s *Store) ListSeepageReadings(ctx context.Context, pointID int64) ([]*model.SeepageReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, point_id, flow_rate, turbidity, ph, temperature, recorded_at, created_at
		FROM seepage_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 100`, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list seepage readings: %v", err)
	}
	defer rows.Close()

	var readings []*model.SeepageReading
	for rows.Next() {
		var r model.SeepageReading
		var recordedAt sql.NullTime
		if err := rows.Scan(&r.ID, &r.PointID, &r.FlowRate, &r.Turbidity, &r.PH,
			&r.Temperature, &recordedAt, &r.CreatedAt); err != nil {
			return nil, nil
		}
		r.RecordedAt = nullTime(recordedAt)
		readings = append(readings, &r)
	}

	return readings, nil
}

// ListAllSeepageReadings 列出所有渗流读数
func (s *Store) ListAllSeepageReadings(ctx context.Context, limit int) ([]*model.SeepageReading, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, point_id, flow_rate, turbidity, ph, temperature, recorded_at, created_at
		FROM seepage_readings ORDER BY recorded_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list all seepage readings: %v", err)
	}
	defer rows.Close()

	var readings []*model.SeepageReading
	for rows.Next() {
		var r model.SeepageReading
		var recordedAt sql.NullTime
		if err := rows.Scan(&r.ID, &r.PointID, &r.FlowRate, &r.Turbidity, &r.PH,
			&r.Temperature, &recordedAt, &r.CreatedAt); err != nil {
			return nil, nil
		}
		r.RecordedAt = nullTime(recordedAt)
		readings = append(readings, &r)
	}

	return readings, nil
}

// CreateDisplacementReading 创建位移读数
func (s *Store) CreateDisplacementReading(ctx context.Context, r *model.DisplacementReading) (*model.DisplacementReading, error) {
	r.CreatedAt = time.Now()
	if r.RecordedAt.IsZero() {
		r.RecordedAt = time.Now()
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO displacement_readings (point_id, horizontal_displacement,
			vertical_displacement, cumulative_horizontal, cumulative_vertical,
			recorded_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		r.PointID, r.HorizontalDisplacement, r.VerticalDisplacement,
		r.CumulativeHorizontal, r.CumulativeVertical, r.RecordedAt, r.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create displacement reading: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get displacement reading id: %v", err)
	}
	r.ID = id

	_ = s.UpdateLastReading(ctx, r.PointID, r.RecordedAt)

	return r, nil
}

// ListDisplacementReadings 列出位移读数
func (s *Store) ListDisplacementReadings(ctx context.Context, pointID int64) ([]*model.DisplacementReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, point_id, horizontal_displacement, vertical_displacement,
			cumulative_horizontal, cumulative_vertical, recorded_at, created_at
		FROM displacement_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 100`, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list displacement readings: %v", err)
	}
	defer rows.Close()

	var readings []*model.DisplacementReading
	for rows.Next() {
		var r model.DisplacementReading
		var recordedAt sql.NullTime
		if err := rows.Scan(&r.ID, &r.PointID, &r.HorizontalDisplacement,
			&r.VerticalDisplacement, &r.CumulativeHorizontal, &r.CumulativeVertical,
			&recordedAt, &r.CreatedAt); err != nil {
			return nil, nil
		}
		r.RecordedAt = nullTime(recordedAt)
		readings = append(readings, &r)
	}

	return readings, nil
}

// CreatePorePressureReading 创建孔隙水压力读数
func (s *Store) CreatePorePressureReading(ctx context.Context, r *model.PorePressureReading) (*model.PorePressureReading, error) {
	r.CreatedAt = time.Now()
	if r.RecordedAt.IsZero() {
		r.RecordedAt = time.Now()
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO pore_pressure_readings (point_id, pressure, depth, water_level,
			recorded_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		r.PointID, r.Pressure, r.Depth, r.WaterLevel, r.RecordedAt, r.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pore pressure reading: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get pore pressure reading id: %v", err)
	}
	r.ID = id

	_ = s.UpdateLastReading(ctx, r.PointID, r.RecordedAt)

	return r, nil
}

// ListPorePressureReadings 列出孔隙水压力读数
func (s *Store) ListPorePressureReadings(ctx context.Context, pointID int64) ([]*model.PorePressureReading, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, point_id, pressure, depth, water_level, recorded_at, created_at
		FROM pore_pressure_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 100`, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list pore pressure readings: %v", err)
	}
	defer rows.Close()

	var readings []*model.PorePressureReading
	for rows.Next() {
		var r model.PorePressureReading
		var recordedAt sql.NullTime
		if err := rows.Scan(&r.ID, &r.PointID, &r.Pressure, &r.Depth, &r.WaterLevel,
			&recordedAt, &r.CreatedAt); err != nil {
			return nil, nil
		}
		r.RecordedAt = nullTime(recordedAt)
		readings = append(readings, &r)
	}

	return readings, nil
}

// BatchCreateReadings 批量创建读数
// bug-007: 遇到错误时 continue 而非返回错误
func (s *Store) BatchCreateReadings(ctx context.Context, batch *model.BatchReadingInput) ([]int64, error) {
	var ids []int64

	for _, input := range batch.SeepageReadings {
		reading := input.ToSeepageReading()
		created, err := s.CreateSeepageReading(ctx, reading)
		if err != nil {
			// bug-007: continue 而非 return
			continue
		}
		ids = append(ids, created.ID)
	}

	for _, input := range batch.DisplacementReadings {
		reading := input.ToDisplacementReading()
		created, err := s.CreateDisplacementReading(ctx, reading)
		if err != nil {
			continue
		}
		ids = append(ids, created.ID)
	}

	for _, input := range batch.PorePressureReadings {
		reading := input.ToPorePressureReading()
		created, err := s.CreatePorePressureReading(ctx, reading)
		if err != nil {
			continue
		}
		ids = append(ids, created.ID)
	}

	return ids, nil
}

// GetLatestSeepageReading 获取最新的渗流读数
func (s *Store) GetLatestSeepageReading(ctx context.Context, pointID int64) (*model.SeepageReading, error) {
	var r model.SeepageReading
	var recordedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, point_id, flow_rate, turbidity, ph, temperature, recorded_at, created_at
		FROM seepage_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 1`, pointID,
	).Scan(&r.ID, &r.PointID, &r.FlowRate, &r.Turbidity, &r.PH, &r.Temperature,
		&recordedAt, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest seepage reading: %v", err)
	}
	r.RecordedAt = nullTime(recordedAt)
	return &r, nil
}

// GetLatestDisplacementReading 获取最新的位移读数
func (s *Store) GetLatestDisplacementReading(ctx context.Context, pointID int64) (*model.DisplacementReading, error) {
	var r model.DisplacementReading
	var recordedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, point_id, horizontal_displacement, vertical_displacement,
			cumulative_horizontal, cumulative_vertical, recorded_at, created_at
		FROM displacement_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 1`, pointID,
	).Scan(&r.ID, &r.PointID, &r.HorizontalDisplacement, &r.VerticalDisplacement,
		&r.CumulativeHorizontal, &r.CumulativeVertical, &recordedAt, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest displacement reading: %v", err)
	}
	r.RecordedAt = nullTime(recordedAt)
	return &r, nil
}

// GetLatestPorePressureReading 获取最新的孔隙水压力读数
func (s *Store) GetLatestPorePressureReading(ctx context.Context, pointID int64) (*model.PorePressureReading, error) {
	var r model.PorePressureReading
	var recordedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, point_id, pressure, depth, water_level, recorded_at, created_at
		FROM pore_pressure_readings WHERE point_id = ? ORDER BY recorded_at DESC LIMIT 1`, pointID,
	).Scan(&r.ID, &r.PointID, &r.Pressure, &r.Depth, &r.WaterLevel,
		&recordedAt, &r.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest pore pressure reading: %v", err)
	}
	r.RecordedAt = nullTime(recordedAt)
	return &r, nil
}
