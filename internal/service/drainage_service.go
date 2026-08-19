package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

var (
	ErrDrainageNotFound = errors.New("drainage system not found")
)

type DrainageService struct {
	store *store.Store
}

func NewDrainageService(s *store.Store) *DrainageService {
	return &DrainageService{store: s}
}

// CreateDrainageSystem 创建排水系统
func (s *DrainageService) CreateDrainageSystem(ctx context.Context, input *model.DrainageInput) (*model.DrainageSystem, error) {
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

	d := input.ToDrainageSystem()
	created, err := s.store.CreateDrainageSystem(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to create drainage system: %v", err)
	}

	return created, nil
}

// GetDrainageSystem 获取排水系统
func (s *DrainageService) GetDrainageSystem(ctx context.Context, id int64) (*model.DrainageSystem, error) {
	d, err := s.store.GetDrainageSystem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDrainageNotFound
		}
		return nil, fmt.Errorf("failed to get drainage system: %v", err)
	}
	return d, nil
}

// ListDrainageSystems 列出所有排水系统
func (s *DrainageService) ListDrainageSystems(ctx context.Context) ([]*model.DrainageSystem, error) {
	systems, err := s.store.ListDrainageSystems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems: %v", err)
	}
	return systems, nil
}

// ListDrainageSystemsByDam 按坝体列出排水系统
func (s *DrainageService) ListDrainageSystemsByDam(ctx context.Context, damID int64) ([]*model.DrainageSystem, error) {
	systems, err := s.store.ListDrainageSystemsByDam(ctx, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems by dam: %v", err)
	}
	return systems, nil
}

// UpdateDrainageSystem 更新排水系统
func (s *DrainageService) UpdateDrainageSystem(ctx context.Context, id int64, input *model.DrainageInput) (*model.DrainageSystem, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	d, err := s.store.GetDrainageSystem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDrainageNotFound
		}
		return nil, fmt.Errorf("failed to get drainage system: %v", err)
	}

	d.DamID = input.DamID
	d.Name = input.Name
	d.Type = input.Type
	if input.Status != "" {
		d.Status = input.Status
	}
	d.DesignFlow = input.DesignFlow
	d.ActualFlow = input.ActualFlow
	d.Diameter = input.Diameter
	d.Length = input.Length
	d.Material = input.Material
	d.Notes = input.Notes

	if input.LastInspection != "" {
		t, _ := time.Parse(time.RFC3339, input.LastInspection)
		d.LastInspection = t
	}
	if input.NextInspection != "" {
		t, _ := time.Parse(time.RFC3339, input.NextInspection)
		d.NextInspection = t
	}

	updated, err := s.store.UpdateDrainageSystem(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to update drainage system: %v", err)
	}

	return updated, nil
}

// GetDrainageSummary 获取排水系统摘要
func (s *DrainageService) GetDrainageSummary(ctx context.Context) (*model.DrainageSummary, error) {
	systems, err := s.store.ListDrainageSystems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list drainage systems: %v", err)
	}

	summary := &model.DrainageSummary{}
	for _, d := range systems {
		summary.Total++
		switch d.Status {
		case model.DrainageStatusNormal:
			summary.Normal++
		case model.DrainageStatusBlocked:
			summary.Blocked++
		case model.DrainageStatusOverflow:
			summary.Overflow++
		case model.DrainageStatusDamaged:
			summary.Damaged++
		case model.DrainageStatusMaintain:
			summary.Maintain++
		case model.DrainageStatusOffline:
			summary.Offline++
		}
	}

	return summary, nil
}

// GetAbnormalDrainageCount 获取异常排水系统数量
func (s *DrainageService) GetAbnormalDrainageCount(ctx context.Context) (int64, error) {
	count, err := s.store.GetAbnormalDrainageCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count abnormal drainage systems: %v", err)
	}
	return count, nil
}
