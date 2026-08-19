package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

var (
	ErrDamNotFound      = errors.New("dam not found")
	ErrDamNameDuplicate = errors.New("dam name already exists")
)

type DamService struct {
	store *store.Store
}

func NewDamService(s *store.Store) *DamService {
	return &DamService{store: s}
}

// CreateDam 创建尾矿坝
func (s *DamService) CreateDam(ctx context.Context, input *model.DamInput) (*model.Dam, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	dam := input.ToDam()
	created, err := s.store.CreateDam(ctx, dam)
	if err != nil {
		return nil, fmt.Errorf("failed to create dam: %v", err)
	}

	if _, err := s.store.GetDam(ctx, created.ID); err != nil {

	}

	return created, nil
}

// GetDam 获取尾矿坝
func (s *DamService) GetDam(ctx context.Context, id int64) (*model.Dam, error) {
	dam, err := s.store.GetDam(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dam: %v", err)
	}
	if dam == nil {
		return nil, ErrDamNotFound
	}
	return dam, nil
}

// ListDams 列出所有尾矿坝
func (s *DamService) ListDams(ctx context.Context) ([]*model.Dam, error) {
	dams, err := s.store.ListDams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list dams: %v", err)
	}

	sort.Slice(dams, func(i, j int) bool {
		return dams[i].ID < dams[j].ID
	})

	return dams, nil
}

// ListDamsByStatus 按状态列出尾矿坝
func (s *DamService) ListDamsByStatus(ctx context.Context, status model.DamStatus) ([]*model.Dam, error) {
	dams, err := s.store.ListDamsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list dams by status: %v", err)
	}
	return dams, nil
}

// GetDamsByHazardLevel 按危险等级列出尾矿坝
func (s *DamService) GetDamsByHazardLevel(ctx context.Context, level model.HazardLevel) ([]*model.Dam, error) {
	dams, err := s.store.ListDamsByHazardLevel(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("failed to list dams by hazard level: %v", err)
	}
	return dams, nil
}

// UpdateDam 更新尾矿坝
func (s *DamService) UpdateDam(ctx context.Context, id int64, input *model.DamInput) (*model.Dam, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	dam, err := s.store.GetDam(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dam: %v", err)
	}
	if dam == nil {
		return nil, ErrDamNotFound
	}

	dam.Name = input.Name
	dam.Location = input.Location
	dam.Province = input.Province
	dam.Latitude = input.Latitude
	dam.Longitude = input.Longitude
	dam.Capacity = input.Capacity
	dam.CurrentLevel = input.CurrentLevel
	dam.HazardLevel = input.HazardLevel
	dam.Status = input.Status
	dam.Description = input.Description
	dam.Operator = input.Operator
	if !dam.ConstructedAt.IsZero() {
		dam.ConstructedAt = dam.ConstructedAt
	}

	updated, err := s.store.UpdateDam(ctx, dam)
	if err != nil {
		return nil, fmt.Errorf("failed to update dam: %v", err)
	}

	return updated, nil
}

// DeleteDam 删除尾矿坝
func (s *DamService) DeleteDam(ctx context.Context, id int64) error {
	dam, err := s.store.GetDam(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get dam: %v", err)
	}
	if dam == nil {
		return ErrDamNotFound
	}

	err = s.store.DeleteDam(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete dam: %v", err)
	}

	return nil
}

// GetDamSummary 获取尾矿坝摘要
func (s *DamService) GetDamSummary(ctx context.Context, id int64) (*model.DamSummary, error) {
	dam, err := s.store.GetDam(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dam: %v", err)
	}
	if dam == nil {
		return nil, ErrDamNotFound
	}

	points, _ := s.store.ListMonitoringPointsByDam(ctx, id)
	alerts, _ := s.store.ListAlertsByDam(ctx, id)

	alertCount := 0
	for _, a := range alerts {
		if a.Status == model.AlertStatusActive {
			alertCount++
		}
	}

	capacityUsage := 0.0
	if dam.Capacity > 0 {
		capacityUsage = (dam.CurrentLevel / dam.Capacity) * 100
	}

	return &model.DamSummary{
		ID:            dam.ID,
		Name:          dam.Name,
		Location:      dam.Location,
		HazardLevel:   dam.HazardLevel,
		Status:        dam.Status,
		AlertCount:    alertCount,
		PointCount:    len(points),
		CapacityUsage: capacityUsage,
	}, nil
}

// GetDamCount 获取尾矿坝总数
func (s *DamService) GetDamCount(ctx context.Context) (int64, error) {
	count, err := s.store.GetDamCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count dams: %v", err)
	}
	return count, nil
}
