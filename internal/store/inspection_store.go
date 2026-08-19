package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tailings-dam/internal/model"
)

// CreateInspection 创建巡检
func (s *Store) CreateInspection(ctx context.Context, insp *model.Inspection) (*model.Inspection, error) {
	now := time.Now()
	insp.CreatedAt = now
	insp.UpdatedAt = now
	if insp.Status == "" {
		insp.Status = model.InspectionStatusPending
	}
	if insp.Priority == "" {
		insp.Priority = model.InspectionPriorityNormal
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO inspections (dam_id, inspector, title, scheduled_date, completed_date,
			findings, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		insp.DamID, insp.Inspector, insp.Title, insp.ScheduledDate,
		nullableTime(insp.CompletedDate), insp.Findings,
		string(insp.Status), string(insp.Priority), insp.CreatedAt, insp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create inspection: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection id: %v", err)
	}
	insp.ID = id

	s.mu.RLock()
	s.inspectionCache[id] = insp
	s.mu.RUnlock()

	return insp, nil
}

// GetInspection 根据 ID 获取巡检
func (s *Store) GetInspection(ctx context.Context, id int64) (*model.Inspection, error) {
	s.mu.RLock()
	if insp, ok := s.inspectionCache[id]; ok {
		s.mu.RUnlock()
		return insp, nil
	}
	s.mu.RUnlock()

	var insp model.Inspection
	var completedDate sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT id, dam_id, inspector, title, scheduled_date, completed_date,
			findings, status, priority, created_at, updated_at
		FROM inspections WHERE id = ?`, id,
	).Scan(&insp.ID, &insp.DamID, &insp.Inspector, &insp.Title, &insp.ScheduledDate,
		&completedDate, &insp.Findings, &insp.Status, &insp.Priority,
		&insp.CreatedAt, &insp.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get inspection: %v", err)
	}
	insp.CompletedDate = nullTime(completedDate)

	s.mu.RLock()
	s.inspectionCache[id] = &insp
	s.mu.RUnlock()

	return &insp, nil
}

// ListInspections 列出所有巡检
func (s *Store) ListInspections(ctx context.Context) ([]*model.Inspection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, inspector, title, scheduled_date, completed_date,
			findings, status, priority, created_at, updated_at
		FROM inspections ORDER BY scheduled_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections: %v", err)
	}
	defer rows.Close()

	return scanInspections(rows)
}

// ListInspectionsByDam 按坝体 ID 列出巡检
func (s *Store) ListInspectionsByDam(ctx context.Context, damID int64) ([]*model.Inspection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, inspector, title, scheduled_date, completed_date,
			findings, status, priority, created_at, updated_at
		FROM inspections WHERE dam_id = ? ORDER BY scheduled_date DESC`, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections by dam: %v", err)
	}
	defer rows.Close()

	return scanInspections(rows)
}

// ListInspectionsByStatus 按状态列出巡检
func (s *Store) ListInspectionsByStatus(ctx context.Context, status model.InspectionStatus) ([]*model.Inspection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, dam_id, inspector, title, scheduled_date, completed_date,
			findings, status, priority, created_at, updated_at
		FROM inspections WHERE status = ? ORDER BY scheduled_date DESC`, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections by status: %v", err)
	}
	defer rows.Close()

	return scanInspections(rows)
}

// UpdateInspection 更新巡检
func (s *Store) UpdateInspection(ctx context.Context, insp *model.Inspection) (*model.Inspection, error) {
	insp.UpdatedAt = time.Now()

	_, err := s.db.ExecContext(ctx,
		`UPDATE inspections SET dam_id = ?, inspector = ?, title = ?, scheduled_date = ?,
			completed_date = ?, findings = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ?`,
		insp.DamID, insp.Inspector, insp.Title, insp.ScheduledDate,
		nullableTime(insp.CompletedDate), insp.Findings,
		string(insp.Status), string(insp.Priority), insp.UpdatedAt, insp.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update inspection: %v", err)
	}

	s.mu.RLock()
	s.inspectionCache[insp.ID] = insp
	s.mu.RUnlock()

	return insp, nil
}

// CompleteInspection 完成巡检
func (s *Store) CompleteInspection(ctx context.Context, id int64, findings string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx,
		`UPDATE inspections SET status = ?, completed_date = ?, findings = ?, updated_at = ?
		WHERE id = ?`,
		string(model.InspectionStatusCompleted), now, findings, now, id)
	if err != nil {
		return fmt.Errorf("failed to complete inspection: %v", err)
	}

	s.mu.RLock()
	if insp, ok := s.inspectionCache[id]; ok {
		insp.Status = model.InspectionStatusCompleted
		insp.CompletedDate = now
		insp.Findings = findings
		insp.UpdatedAt = now
	}
	s.mu.RUnlock()

	return nil
}

// scanInspections 扫描巡检行集
func scanInspections(rows *sql.Rows) ([]*model.Inspection, error) {
	var inspections []*model.Inspection
	for rows.Next() {
		var insp model.Inspection
		var completedDate sql.NullTime
		if err := rows.Scan(&insp.ID, &insp.DamID, &insp.Inspector, &insp.Title,
			&insp.ScheduledDate, &completedDate, &insp.Findings, &insp.Status,
			&insp.Priority, &insp.CreatedAt, &insp.UpdatedAt); err != nil {
			return nil, nil
		}
		insp.CompletedDate = nullTime(completedDate)
		inspections = append(inspections, &insp)
	}
	return inspections, nil
}

// GetInspectionCount 获取巡检总数
func (s *Store) GetInspectionCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM inspections").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count inspections: %v", err)
	}
	return count, nil
}

// GetPendingInspectionCount 获取待巡检数量
func (s *Store) GetPendingInspectionCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM inspections WHERE status = ?",
		string(model.InspectionStatusPending)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending inspections: %v", err)
	}
	return count, nil
}

// BatchCreateInspections 批量创建巡检
func (s *Store) BatchCreateInspections(ctx context.Context, inspections []*model.Inspection) ([]int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // 在 Commit 成功后调用为 no-op

	var ids []int64
	for _, insp := range inspections {
		now := time.Now()
		insp.CreatedAt = now
		insp.UpdatedAt = now
		if insp.Status == "" {
			insp.Status = model.InspectionStatusPending
		}
		if insp.Priority == "" {
			insp.Priority = model.InspectionPriorityNormal
		}
		result, err := tx.ExecContext(ctx,
			`INSERT INTO inspections (dam_id, inspector, title, scheduled_date, completed_date,
				findings, status, priority, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			insp.DamID, insp.Inspector, insp.Title, insp.ScheduledDate,
			nullableTime(insp.CompletedDate), insp.Findings,
			string(insp.Status), string(insp.Priority), insp.CreatedAt, insp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert inspection: %v", err)
		}
		id, _ := result.LastInsertId()
		insp.ID = id
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	return ids, nil
}
