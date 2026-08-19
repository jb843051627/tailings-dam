package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateDrainageSystem 创建排水系统
func (s *Store) CreateDrainageSystem(ctx context.Context, d *model.DrainageSystem) (*model.DrainageSystem, error) {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Status == "" {
		d.Status = model.DrainageStatusNormal
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO drainage_systems (dam_id, name, type, status, design_flow,
			actual_flow, diameter, length, material, last_inspection, next_inspection,
			notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.DamID, d.Name, string(d.Type), string(d.Status), d.DesignFlow,
		d.ActualFlow, d.Diameter, d.Length, d.Material,
		nullableTime(d.LastInspection), nullableTime(d.NextInspection),
		d.Notes, d.CreatedAt, d.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create drainage system: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get drainage system id: %v", err)
	}
	d.ID = id

	s.mu.RLock()
	s.drainageCache[id] = d
	s.mu.RUnlock()

	return d, nil
}

// GetDrainageSystem 根据 ID 获取排水系统
func (s *Store) GetDrainageSystem(ctx context.Context, id int64) (*model.DrainageSystem, error) {
	s.mu.RLock()
	if d, ok := s.drainageCache[id]; ok {
		s.mu.RUnlock()
		return d, nil
	}
	s.mu.RUnlock()

	var d model.DrainageSystem
	var prevInspection, nextInspection sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, dam_id, name, type, status, design_flow, actual_flow, diameter,
			length, material, last_inspection, next_inspection, notes, created_at, updated_at
		FROM drainage_systems WHERE id = ?`, id,
	).Scan(&d.ID, &d.DamID, &d.Name, &d.Type, &d.Status, &d.DesignFlow,
		&d.ActualFlow, &d.Diameter, &d.Length, &d.Material,
		&prevInspection, &nextInspection, &d.Notes, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("drainage system not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get drainage system: %v", err)
	}
	d.LastInspection = nullTime(prevInspection)
	d.NextInspection = nullTime(nextInspection)

	s.mu.RLock()
	s.drainageCache[id] = &d
	s.mu.RUnlock()

	return &d, nil
}

// ListDrainageSystems 列出所有排水系统
func (s *Store) ListDrainageSystems(ctx context.Context) ([]*model.DrainageSystem, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, type, status, design_flow, actual_flow, diameter,
			length, material, last_inspection, next_inspection, notes, created_at, updated_at
		FROM drainage_systems ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems: %v", err)
	}
	defer rows.Close()

	return scanDrainageSystems(rows)
}

// ListDrainageSystemsByDam 按坝体 ID 列出排水系统
func (s *Store) ListDrainageSystemsByDam(ctx context.Context, damID int64) ([]*model.DrainageSystem, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, type, status, design_flow, actual_flow, diameter,
			length, material, last_inspection, next_inspection, notes, created_at, updated_at
		FROM drainage_systems WHERE dam_id = ? ORDER BY id`, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems by dam: %v", err)
	}
	defer rows.Close()

	return scanDrainageSystems(rows)
}

// ListDrainageSystemsByStatus 按状态列出排水系统
func (s *Store) ListDrainageSystemsByStatus(ctx context.Context, status model.DrainageStatus) ([]*model.DrainageSystem, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, name, type, status, design_flow, actual_flow, diameter,
			length, material, last_inspection, next_inspection, notes, created_at, updated_at
		FROM drainage_systems WHERE status = ? ORDER BY id`, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems by status: %v", err)
	}
	defer rows.Close()

	return scanDrainageSystems(rows)
}

// UpdateDrainageSystem 更新排水系统
func (s *Store) UpdateDrainageSystem(ctx context.Context, d *model.DrainageSystem) (*model.DrainageSystem, error) {
	d.UpdatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`UPDATE drainage_systems SET dam_id = ?, name = ?, type = ?, status = ?,
			design_flow = ?, actual_flow = ?, diameter = ?, length = ?, material = ?,
			last_inspection = ?, next_inspection = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		d.DamID, d.Name, string(d.Type), string(d.Status), d.DesignFlow,
		d.ActualFlow, d.Diameter, d.Length, d.Material,
		nullableTime(d.LastInspection), nullableTime(d.NextInspection),
		d.Notes, d.UpdatedAt, d.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update drainage system: %v", err)
	}

	s.mu.RLock()
	s.drainageCache[d.ID] = d
	s.mu.RUnlock()

	return d, nil
}

// DeleteDrainageSystem 删除排水系统
func (s *Store) DeleteDrainageSystem(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM drainage_systems WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete drainage system: %v", err)
	}

	s.mu.RLock()
	delete(s.drainageCache, id)
	s.mu.RUnlock()

	return nil
}

// scanDrainageSystems 扫描排水系统行集
func scanDrainageSystems(rows *sql.Rows) ([]*model.DrainageSystem, error) {
	var systems []*model.DrainageSystem
	for rows.Next() {
		var d model.DrainageSystem
		var prevInspection, nextInspection sql.NullTime
		if err := rows.Scan(&d.ID, &d.DamID, &d.Name, &d.Type, &d.Status,
			&d.DesignFlow, &d.ActualFlow, &d.Diameter, &d.Length, &d.Material,
			&prevInspection, &nextInspection, &d.Notes, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, nil
		}
		d.LastInspection = nullTime(prevInspection)
		d.NextInspection = nullTime(nextInspection)
		systems = append(systems, &d)
	}
	return systems, nil
}

// GetDrainageSystemCount 获取排水系统总数
func (s *Store) GetDrainageSystemCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM drainage_systems").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count drainage systems: %v", err)
	}
	return count, nil
}

// GetAbnormalDrainageCount 获取异常排水系统数量
func (s *Store) GetAbnormalDrainageCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM drainage_systems WHERE status != ?`,
		string(model.DrainageStatusNormal)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count abnormal drainage systems: %v", err)
	}
	return count, nil
}
