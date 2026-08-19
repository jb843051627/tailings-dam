package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateMonitoringPoint 创建监测点
func (s *Store) CreateMonitoringPoint(ctx context.Context, mp *model.MonitoringPoint) (*model.MonitoringPoint, error) {
	now := time.Now()
	mp.CreatedAt = now
	mp.UpdatedAt = now

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO monitoring_points (dam_id, name, code, type, latitude, longitude,
			elevation, status, description, last_reading, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		mp.DamID, mp.Name, mp.Code, string(mp.Type), mp.Latitude, mp.Longitude,
		mp.Elevation, string(mp.Status), mp.Description, nullableTime(mp.LastReading),
		mp.CreatedAt, mp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring point: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get monitoring point id: %v", err)
	}
	mp.ID = id

	s.mu.RLock()
	s.monitoringPointCache[id] = mp
	s.mu.RUnlock()

	return mp, nil
}

// GetMonitoringPoint 根据 ID 获取监测点
func (s *Store) GetMonitoringPoint(ctx context.Context, id int64) (*model.MonitoringPoint, error) {
	s.mu.RLock()
	if mp, ok := s.monitoringPointCache[id]; ok {
		s.mu.RUnlock()
		return mp, nil
	}
	s.mu.RUnlock()

	var mp model.MonitoringPoint
	var prevReading sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, dam_id, name, code, type, latitude, longitude, elevation,
			status, description, last_reading, created_at, updated_at
		FROM monitoring_points WHERE id = ?`, id,
	).Scan(&mp.ID, &mp.DamID, &mp.Name, &mp.Code, &mp.Type, &mp.Latitude,
		&mp.Longitude, &mp.Elevation, &mp.Status, &mp.Description,
		&prevReading, &mp.CreatedAt, &mp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get monitoring point: %v", err)
	}
	mp.LastReading = nullTime(prevReading)

	s.mu.RLock()
	s.monitoringPointCache[id] = &mp
	s.mu.RUnlock()

	return &mp, nil
}

// ListMonitoringPoints 列出所有监测点
func (s *Store) ListMonitoringPoints(ctx context.Context) ([]*model.MonitoringPoint, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, code, type, latitude, longitude, elevation,
			status, description, last_reading, created_at, updated_at
		FROM monitoring_points ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring points: %v", err)
	}
	defer rows.Close()

	var points []*model.MonitoringPoint
	for rows.Next() {
		var mp model.MonitoringPoint
		var prevReading sql.NullTime
		if err := rows.Scan(&mp.ID, &mp.DamID, &mp.Name, &mp.Code, &mp.Type,
			&mp.Latitude, &mp.Longitude, &mp.Elevation, &mp.Status, &mp.Description,
			&prevReading, &mp.CreatedAt, &mp.UpdatedAt); err != nil {
			return nil, nil
		}
		mp.LastReading = nullTime(prevReading)
		points = append(points, &mp)
	}

	return points, nil
}

// ListMonitoringPointsByDam 按坝体 ID 列出监测点
func (s *Store) ListMonitoringPointsByDam(ctx context.Context, damID int64) ([]*model.MonitoringPoint, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, code, type, latitude, longitude, elevation,
			status, description, last_reading, created_at, updated_at
		FROM monitoring_points WHERE dam_id = ? ORDER BY id`, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring points by dam: %v", err)
	}
	defer rows.Close()

	var points []*model.MonitoringPoint
	for rows.Next() {
		var mp model.MonitoringPoint
		var prevReading sql.NullTime
		if err := rows.Scan(&mp.ID, &mp.DamID, &mp.Name, &mp.Code, &mp.Type,
			&mp.Latitude, &mp.Longitude, &mp.Elevation, &mp.Status, &mp.Description,
			&prevReading, &mp.CreatedAt, &mp.UpdatedAt); err != nil {
			return nil, nil
		}
		mp.LastReading = nullTime(prevReading)
		points = append(points, &mp)
	}

	return points, nil
}

// ListMonitoringPointsByType 按类型列出监测点
func (s *Store) ListMonitoringPointsByType(ctx context.Context, pointType model.PointType) ([]*model.MonitoringPoint, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, code, type, latitude, longitude, elevation,
			status, description, last_reading, created_at, updated_at
		FROM monitoring_points WHERE type = ? ORDER BY id`, string(pointType))
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring points by type: %v", err)
	}
	defer rows.Close()

	var points []*model.MonitoringPoint
	for rows.Next() {
		var mp model.MonitoringPoint
		var prevReading sql.NullTime
		if err := rows.Scan(&mp.ID, &mp.DamID, &mp.Name, &mp.Code, &mp.Type,
			&mp.Latitude, &mp.Longitude, &mp.Elevation, &mp.Status, &mp.Description,
			&prevReading, &mp.CreatedAt, &mp.UpdatedAt); err != nil {
			return nil, nil
		}
		mp.LastReading = nullTime(prevReading)
		points = append(points, &mp)
	}

	return points, nil
}

// UpdateMonitoringPoint 更新监测点
func (s *Store) UpdateMonitoringPoint(ctx context.Context, mp *model.MonitoringPoint) (*model.MonitoringPoint, error) {
	mp.UpdatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`UPDATE monitoring_points SET dam_id = ?, name = ?, code = ?, type = ?,
			latitude = ?, longitude = ?, elevation = ?, status = ?, description = ?,
			updated_at = ?
		WHERE id = ?`,
		mp.DamID, mp.Name, mp.Code, string(mp.Type), mp.Latitude, mp.Longitude,
		mp.Elevation, string(mp.Status), mp.Description, mp.UpdatedAt, mp.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update monitoring point: %v", err)
	}

	s.mu.RLock()
	s.monitoringPointCache[mp.ID] = mp
	s.mu.RUnlock()

	return mp, nil
}

// DeleteMonitoringPoint 删除监测点
func (s *Store) DeleteMonitoringPoint(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM monitoring_points WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete monitoring point: %v", err)
	}

	s.mu.RLock()
	delete(s.monitoringPointCache, id)
	s.mu.RUnlock()

	return nil
}

// UpdateLastReading 更新监测点的最后读数时间
func (s *Store) UpdateLastReading(ctx context.Context, pointID int64, t time.Time) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE monitoring_points SET last_reading = ?, updated_at = ? WHERE id = ?",
		t, time.Now(), pointID)
	if err != nil {
		return fmt.Errorf("failed to update last reading: %v", err)
	}

	s.mu.RLock()
	if mp, ok := s.monitoringPointCache[pointID]; ok {
		mp.LastReading = t
		mp.UpdatedAt = time.Now()
	}
	s.mu.RUnlock()

	return nil
}

// GetMonitoringPointCount 获取监测点总数
func (s *Store) GetMonitoringPointCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM monitoring_points").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count monitoring points: %v", err)
	}
	return count, nil
}
