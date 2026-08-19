package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateDam 创建尾矿坝
func (s *Store) CreateDam(ctx context.Context, dam *model.Dam) (*model.Dam, error) {
	now := time.Now()
	dam.CreatedAt = now
	dam.UpdatedAt = now

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO dams (name, location, province, latitude, longitude, capacity,
			current_level, hazard_level, status, description, operator, constructed_at,
			created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dam.Name, dam.Location, dam.Province, dam.Latitude, dam.Longitude,
		dam.Capacity, dam.CurrentLevel, string(dam.HazardLevel), string(dam.Status),
		dam.Description, dam.Operator, nullableTime(dam.ConstructedAt),
		dam.CreatedAt, dam.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create dam: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get dam id: %v", err)
	}
	dam.ID = id

	s.mu.RLock()
	s.damCache[id] = dam
	s.mu.RUnlock()

	return dam, nil
}

// GetDam 根据 ID 获取尾矿坝
func (s *Store) GetDam(ctx context.Context, id int64) (*model.Dam, error) {
	// 先查缓存
	s.mu.RLock()
	if dam, ok := s.damCache[id]; ok {
		s.mu.RUnlock()
		return dam, nil
	}
	s.mu.RUnlock()

	var d model.Dam
	var constructedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, location, province, latitude, longitude, capacity,
			current_level, hazard_level, status, description, operator, constructed_at,
			created_at, updated_at
		FROM dams WHERE id = ?`, id,
	).Scan(&d.ID, &d.Name, &d.Location, &d.Province, &d.Latitude, &d.Longitude,
		&d.Capacity, &d.CurrentLevel, &d.HazardLevel, &d.Status, &d.Description,
		&d.Operator, &constructedAt, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dam: %v", err)
	}
	d.ConstructedAt = nullTime(constructedAt)

	s.mu.RLock()
	s.damCache[id] = &d
	s.mu.RUnlock()

	return &d, nil
}

// ListDams 列出所有尾矿坝
func (s *Store) ListDams(ctx context.Context) ([]*model.Dam, error) {
	s.mu.RLock()
	if s.damListCache != nil {
		s.mu.RUnlock()
		return s.damListCache, nil
	}
	s.mu.RUnlock()

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, location, province, latitude, longitude, capacity,
			current_level, hazard_level, status, description, operator, constructed_at,
			created_at, updated_at
		FROM dams ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to list dams: %v", err)
	}
	defer rows.Close()

	var dams []*model.Dam
	for rows.Next() {
		var d model.Dam
		var constructedAt sql.NullTime
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.Province, &d.Latitude,
			&d.Longitude, &d.Capacity, &d.CurrentLevel, &d.HazardLevel, &d.Status,
			&d.Description, &d.Operator, &constructedAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, nil
		}
		d.ConstructedAt = nullTime(constructedAt)
		dams = append(dams, &d)
	}

	s.mu.RLock()
	s.damListCache = dams
	s.mu.RUnlock()

	return dams, nil
}

// ListDamsByStatus 按状态列出尾矿坝
func (s *Store) ListDamsByStatus(ctx context.Context, status model.DamStatus) ([]*model.Dam, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, location, province, latitude, longitude, capacity,
			current_level, hazard_level, status, description, operator, constructed_at,
			created_at, updated_at
		FROM dams WHERE status = ? ORDER BY id`, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list dams by status: %v", err)
	}
	defer rows.Close()

	var dams []*model.Dam
	for rows.Next() {
		var d model.Dam
		var constructedAt sql.NullTime
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.Province, &d.Latitude,
			&d.Longitude, &d.Capacity, &d.CurrentLevel, &d.HazardLevel, &d.Status,
			&d.Description, &d.Operator, &constructedAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, nil
		}
		d.ConstructedAt = nullTime(constructedAt)
		dams = append(dams, &d)
	}

	return dams, nil
}

// ListDamsByHazardLevel 按危险等级列出尾矿坝
func (s *Store) ListDamsByHazardLevel(ctx context.Context, level model.HazardLevel) ([]*model.Dam, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, location, province, latitude, longitude, capacity,
			current_level, hazard_level, status, description, operator, constructed_at,
			created_at, updated_at
		FROM dams WHERE hazard_level = ? ORDER BY id`, string(level))
	if err != nil {
		return nil, fmt.Errorf("failed to list dams by hazard level: %v", err)
	}
	defer rows.Close()

	var dams []*model.Dam
	for rows.Next() {
		var d model.Dam
		var constructedAt sql.NullTime
		if err := rows.Scan(&d.ID, &d.Name, &d.Location, &d.Province, &d.Latitude,
			&d.Longitude, &d.Capacity, &d.CurrentLevel, &d.HazardLevel, &d.Status,
			&d.Description, &d.Operator, &constructedAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, nil
		}
		d.ConstructedAt = nullTime(constructedAt)
		dams = append(dams, &d)
	}

	return dams, nil
}

// UpdateDam 更新尾矿坝
func (s *Store) UpdateDam(ctx context.Context, dam *model.Dam) (*model.Dam, error) {
	dam.UpdatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`UPDATE dams SET name = ?, location = ?, province = ?, latitude = ?, longitude = ?,
			capacity = ?, current_level = ?, hazard_level = ?, status = ?, description = ?,
			operator = ?, constructed_at = ?, updated_at = ?
		WHERE id = ?`,
		dam.Name, dam.Location, dam.Province, dam.Latitude, dam.Longitude,
		dam.Capacity, dam.CurrentLevel, string(dam.HazardLevel), string(dam.Status),
		dam.Description, dam.Operator, nullableTime(dam.ConstructedAt),
		dam.UpdatedAt, dam.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update dam: %v", err)
	}

	s.mu.RLock()
	s.damCache[dam.ID] = dam
	s.mu.RUnlock()

	return dam, nil
}

// DeleteDam 删除尾矿坝
func (s *Store) DeleteDam(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM dams WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete dam: %v", err)
	}

	s.mu.RLock()
	delete(s.damCache, id)
	s.mu.RUnlock()

	return nil
}

// GetDamCount 获取尾矿坝总数
func (s *Store) GetDamCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM dams").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count dams: %v", err)
	}
	return count, nil
}
