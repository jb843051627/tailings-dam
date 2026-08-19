package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

var (
	ErrInspectionNotFound = errors.New("inspection not found")
	ErrInspectionCompleted = errors.New("inspection already completed")
)

type InspectionService struct {
	store *store.Store
}

func NewInspectionService(s *store.Store) *InspectionService {
	return &InspectionService{store: s}
}

// CreateInspection 创建巡检
func (s *InspectionService) CreateInspection(ctx context.Context, input *model.InspectionInput) (*model.Inspection, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	// 验证坝体是否存在
	dam, err := s.store.GetDam(ctx, input.DamID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify dam: %v", err)
	}
	if dam == nil {
		return nil, ErrPointDamNotFound
	}

	insp := input.ToInspection()
	created, err := s.store.CreateInspection(ctx, insp)
	if err != nil {
		return nil, fmt.Errorf("failed to create inspection: %v", err)
	}

	return created, nil
}

// GetInspection 获取巡检
func (s *InspectionService) GetInspection(ctx context.Context, id int64) (*model.Inspection, error) {
	insp, err := s.store.GetInspection(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection: %v", err)
	}
	if insp == nil {
		return nil, ErrInspectionNotFound
	}
	return insp, nil
}

// ListInspections 列出所有巡检
func (s *InspectionService) ListInspections(ctx context.Context) ([]*model.Inspection, error) {
	inspections, err := s.store.ListInspections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections: %v", err)
	}

	now := time.Now()
	for _, insp := range inspections {
		if insp.Status == model.InspectionStatusPending {
			if insp.ScheduledDate.Before(now) {
				insp.Status = model.InspectionStatusOverdue
			}
		}
	}

	return inspections, nil
}

// ListInspectionsByDam 按坝体列出巡检
func (s *InspectionService) ListInspectionsByDam(ctx context.Context, damID int64) ([]*model.Inspection, error) {
	inspections, err := s.store.ListInspectionsByDam(ctx, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections by dam: %v", err)
	}
	return inspections, nil
}

// CompleteInspection 完成巡检
func (s *InspectionService) CompleteInspection(ctx context.Context, id int64, findings string) error {
	insp, err := s.store.GetInspection(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get inspection: %v", err)
	}
	if insp == nil {
		return ErrInspectionNotFound
	}

	if insp.Status == model.InspectionStatusCompleted {
		return ErrInspectionCompleted
	}

	err = s.store.CompleteInspection(ctx, id, findings)
	if err != nil {
		return fmt.Errorf("failed to complete inspection: %v", err)
	}

	return nil
}

// GetInspectionSummary 获取巡检摘要
func (s *InspectionService) GetInspectionSummary(ctx context.Context) (*model.InspectionSummary, error) {
	inspections, err := s.store.ListInspections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list inspections: %v", err)
	}

	summary := &model.InspectionSummary{}
	now := time.Now()
	for _, insp := range inspections {
		summary.Total++
		switch insp.Status {
		case model.InspectionStatusPending:
			if insp.ScheduledDate.Before(now) {
				summary.Overdue++
			} else {
				summary.Pending++
			}
		case model.InspectionStatusInProgress:
			summary.InProgress++
		case model.InspectionStatusCompleted:
			summary.Completed++
		case model.InspectionStatusOverdue:
			summary.Overdue++
		}
	}

	return summary, nil
}

// UpdateInspection 更新巡检
func (s *InspectionService) UpdateInspection(ctx context.Context, id int64, input *model.InspectionInput) (*model.Inspection, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	insp, err := s.store.GetInspection(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get inspection: %v", err)
	}
	if insp == nil {
		return nil, ErrInspectionNotFound
	}

	insp.DamID = input.DamID
	insp.Inspector = input.Inspector
	insp.Title = input.Title
	insp.Priority = input.Priority
	if input.ScheduledDate != "" {
		t, _ := time.Parse(time.RFC3339, input.ScheduledDate)
		insp.ScheduledDate = t
	}
	if input.Findings != "" {
		insp.Findings = input.Findings
	}

	updated, err := s.store.UpdateInspection(ctx, insp)
	if err != nil {
		return nil, fmt.Errorf("failed to update inspection: %v", err)
	}

	return updated, nil
}
