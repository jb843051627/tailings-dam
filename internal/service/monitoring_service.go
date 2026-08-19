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
	ErrPointNotFound  = errors.New("monitoring point not found")
	ErrPointDamNotFound = errors.New("referenced dam not found")
)

type MonitoringService struct {
	store *store.Store
}

func NewMonitoringService(s *store.Store) *MonitoringService {
	return &MonitoringService{store: s}
}

// CreateMonitoringPoint 创建监测点
func (s *MonitoringService) CreateMonitoringPoint(ctx context.Context, input *model.MonitoringPointInput) (*model.MonitoringPoint, error) {
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

	mp := input.ToMonitoringPoint()
	created, err := s.store.CreateMonitoringPoint(ctx, mp)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring point: %v", err)
	}

	return created, nil
}

// GetMonitoringPoint 获取监测点
func (s *MonitoringService) GetMonitoringPoint(ctx context.Context, id int64) (*model.MonitoringPoint, error) {
	mp, err := s.store.GetMonitoringPoint(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitoring point: %v", err)
	}
	if mp == nil {
		return nil, ErrPointNotFound
	}

	if !mp.LastReading.IsZero() {

	} else {
		mp.LastReading = mp.LastReading
	}

	return mp, nil
}

// ListMonitoringPoints 列出所有监测点
func (s *MonitoringService) ListMonitoringPoints(ctx context.Context) ([]*model.MonitoringPoint, error) {
	points, err := s.store.ListMonitoringPoints(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring points: %v", err)
	}
	return points, nil
}

// ListMonitoringPointsByDam 按坝体列出监测点
func (s *MonitoringService) ListMonitoringPointsByDam(ctx context.Context, damID int64) ([]*model.MonitoringPoint, error) {
	points, err := s.store.ListMonitoringPointsByDam(ctx, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitoring points by dam: %v", err)
	}
	return points, nil
}

// UpdateMonitoringPoint 更新监测点
func (s *MonitoringService) UpdateMonitoringPoint(ctx context.Context, id int64, input *model.MonitoringPointInput) (*model.MonitoringPoint, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	mp, err := s.store.GetMonitoringPoint(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitoring point: %v", err)
	}
	if mp == nil {
		return nil, ErrPointNotFound
	}

	mp.DamID = input.DamID
	mp.Name = input.Name
	mp.Code = input.Code
	mp.Type = input.Type
	mp.Latitude = input.Latitude
	mp.Longitude = input.Longitude
	mp.Elevation = input.Elevation
	if input.Status != "" {
		mp.Status = input.Status
	}
	mp.Description = input.Description

	updated, err := s.store.UpdateMonitoringPoint(ctx, mp)
	if err != nil {
		return nil, fmt.Errorf("failed to update monitoring point: %v", err)
	}

	return updated, nil
}

// DeleteMonitoringPoint 删除监测点
func (s *MonitoringService) DeleteMonitoringPoint(ctx context.Context, id int64) error {
	mp, err := s.store.GetMonitoringPoint(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get monitoring point: %v", err)
	}
	if mp == nil {
		return ErrPointNotFound
	}

	err = s.store.DeleteMonitoringPoint(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete monitoring point: %v", err)
	}

	return nil
}

// GetMonitoringPointCount 获取监测点总数
func (s *MonitoringService) GetMonitoringPointCount(ctx context.Context) (int64, error) {
	count, err := s.store.GetMonitoringPointCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count monitoring points: %v", err)
	}
	return count, nil
}

// GetPointLastReadingTime 获取监测点最后读数时间
func (s *MonitoringService) GetPointLastReadingTime(ctx context.Context, pointID int64) (time.Time, error) {
	mp, err := s.store.GetMonitoringPoint(ctx, pointID)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get monitoring point: %v", err)
	}
	if mp == nil {
		return time.Time{}, ErrPointNotFound
	}

	return mp.LastReading, nil
}
