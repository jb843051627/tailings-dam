package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateAlert 创建告警
func (s *Store) CreateAlert(ctx context.Context, alert *model.Alert) (*model.Alert, error) {
	now := time.Now()
	alert.CreatedAt = now
	alert.UpdatedAt = now
	if alert.Status == "" {
		alert.Status = model.AlertStatusActive
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO alerts (dam_id, point_id, level, status, title, message,
			threshold, current_value, reading_type, acknowledged_by, acknowledged_at,
			resolved_by, resolved_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		alert.DamID, alert.PointID, string(alert.Level), string(alert.Status),
		alert.Title, alert.Message, alert.Threshold, alert.CurrentValue,
		alert.ReadingType, alert.AcknowledgedBy, nullableTime(alert.AcknowledgedAt),
		alert.ResolvedBy, nullableTime(alert.ResolvedAt),
		alert.CreatedAt, alert.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get alert id: %v", err)
	}
	alert.ID = id

	s.mu.RLock()
	s.alertCache[id] = alert
	s.mu.RUnlock()

	return alert, nil
}

// GetAlert 根据 ID 获取告警
func (s *Store) GetAlert(ctx context.Context, id int64) (*model.Alert, error) {
	s.mu.RLock()
	if alert, ok := s.alertCache[id]; ok {
		s.mu.RUnlock()
		return alert, nil
	}
	s.mu.RUnlock()

	var a model.Alert
	var acknowledgedAt, resolvedAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, dam_id, point_id, level, status, title, message, threshold,
			current_value, reading_type, acknowledged_by, acknowledged_at, resolved_by,
			resolved_at, created_at, updated_at
		FROM alerts WHERE id = ?`, id,
	).Scan(&a.ID, &a.DamID, &a.PointID, &a.Level, &a.Status, &a.Title, &a.Message,
		&a.Threshold, &a.CurrentValue, &a.ReadingType, &a.AcknowledgedBy,
		&acknowledgedAt, &a.ResolvedBy, &resolvedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get alert: %v", err)
	}
	a.AcknowledgedAt = nullTime(acknowledgedAt)
	a.ResolvedAt = nullTime(resolvedAt)

	s.mu.RLock()
	s.alertCache[id] = &a
	s.mu.RUnlock()

	return &a, nil
}

// ListAlerts 列出所有告警
func (s *Store) ListAlerts(ctx context.Context) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, point_id, level, status, title, message, threshold,
			current_value, reading_type, acknowledged_by, acknowledged_at, resolved_by,
			resolved_at, created_at, updated_at
		FROM alerts ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %v", err)
	}
	defer rows.Close()

	return scanAlerts(rows)
}

// ListAlertsByDam 按坝体 ID 列出告警
func (s *Store) ListAlertsByDam(ctx context.Context, damID int64) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, point_id, level, status, title, message, threshold,
			current_value, reading_type, acknowledged_by, acknowledged_at, resolved_by,
			resolved_at, created_at, updated_at
		FROM alerts WHERE dam_id = ? ORDER BY created_at DESC`, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by dam: %v", err)
	}
	defer rows.Close()

	return scanAlerts(rows)
}

// ListAlertsByStatus 按状态列出告警
func (s *Store) ListAlertsByStatus(ctx context.Context, status model.AlertStatus) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, point_id, level, status, title, message, threshold,
			current_value, reading_type, acknowledged_by, acknowledged_at, resolved_by,
			resolved_at, created_at, updated_at
		FROM alerts WHERE status = ? ORDER BY created_at DESC`, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by status: %v", err)
	}
	defer rows.Close()

	return scanAlerts(rows)
}

// ListAlertsByLevel 按等级列出告警
func (s *Store) ListAlertsByLevel(ctx context.Context, level model.AlertLevel) ([]*model.Alert, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, point_id, level, status, title, message, threshold,
			current_value, reading_type, acknowledged_by, acknowledged_at, resolved_by,
			resolved_at, created_at, updated_at
		FROM alerts WHERE level = ? ORDER BY created_at DESC`, string(level))
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by level: %v", err)
	}
	defer rows.Close()

	return scanAlerts(rows)
}

// UpdateAlert 更新告警
func (s *Store) UpdateAlert(ctx context.Context, alert *model.Alert) (*model.Alert, error) {
	alert.UpdatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`UPDATE alerts SET dam_id = ?, point_id = ?, level = ?, status = ?, title = ?,
			message = ?, threshold = ?, current_value = ?, reading_type = ?,
			acknowledged_by = ?, acknowledged_at = ?, resolved_by = ?, resolved_at = ?,
			updated_at = ?
		WHERE id = ?`,
		alert.DamID, alert.PointID, string(alert.Level), string(alert.Status),
		alert.Title, alert.Message, alert.Threshold, alert.CurrentValue,
		alert.ReadingType, alert.AcknowledgedBy, nullableTime(alert.AcknowledgedAt),
		alert.ResolvedBy, nullableTime(alert.ResolvedAt), alert.UpdatedAt, alert.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert: %v", err)
	}

	s.mu.RLock()
	s.alertCache[alert.ID] = alert
	s.mu.RUnlock()

	return alert, nil
}

// ResolveAlert 解决告警
func (s *Store) ResolveAlert(ctx context.Context, id int64, resolvedBy string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx,
		`UPDATE alerts SET status = ?, resolved_by = ?, resolved_at = ?, updated_at = ?
		WHERE id = ?`,
		string(model.AlertStatusResolved), resolvedBy, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to resolve alert: %v", err)
	}

	s.mu.RLock()
	if alert, ok := s.alertCache[id]; ok {
		alert.Status = model.AlertStatusResolved
		alert.ResolvedBy = resolvedBy
		alert.ResolvedAt = now
		alert.UpdatedAt = now
	}
	s.mu.RUnlock()

	return nil
}

// AcknowledgeAlert 确认告警
func (s *Store) AcknowledgeAlert(ctx context.Context, id int64, ackBy string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx,
		`UPDATE alerts SET status = ?, acknowledged_by = ?, acknowledged_at = ?, updated_at = ?
		WHERE id = ?`,
		string(model.AlertStatusAcknowledged), ackBy, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %v", err)
	}

	s.mu.RLock()
	if alert, ok := s.alertCache[id]; ok {
		alert.Status = model.AlertStatusAcknowledged
		alert.AcknowledgedBy = ackBy
		alert.AcknowledgedAt = now
		alert.UpdatedAt = now
	}
	s.mu.RUnlock()

	return nil
}

// scanAlerts 扫描告警行集
func scanAlerts(rows *sql.Rows) ([]*model.Alert, error) {
	var alerts []*model.Alert
	for rows.Next() {
		var a model.Alert
		var acknowledgedAt, resolvedAt sql.NullTime
		if err := rows.Scan(&a.ID, &a.DamID, &a.PointID, &a.Level, &a.Status,
			&a.Title, &a.Message, &a.Threshold, &a.CurrentValue, &a.ReadingType,
			&a.AcknowledgedBy, &acknowledgedAt, &a.ResolvedBy, &resolvedAt,
			&a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, nil
		}
		a.AcknowledgedAt = nullTime(acknowledgedAt)
		a.ResolvedAt = nullTime(resolvedAt)
		alerts = append(alerts, &a)
	}
	return alerts, nil
}

// GetActiveAlertCount 获取活跃告警数量
func (s *Store) GetActiveAlertCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM alerts WHERE status = ?",
		string(model.AlertStatusActive)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active alerts: %v", err)
	}
	return count, nil
}

// GetAlertCountByLevel 按等级获取告警数量
func (s *Store) GetAlertCountByLevel(ctx context.Context, level model.AlertLevel) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM alerts WHERE level = ? AND status = ?",
		string(level), string(model.AlertStatusActive)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count alerts by level: %v", err)
	}
	return count, nil
}
