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
	ErrAlertNotFound       = errors.New("alert not found")
	ErrAlertAlreadyResolved = errors.New("alert already resolved")
)

type AlertService struct {
	store *store.Store
}

func NewAlertService(s *store.Store) *AlertService {
	return &AlertService{store: s}
}

// CreateAlert 创建告警
func (s *AlertService) CreateAlert(ctx context.Context, input *model.AlertInput) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	alert := input.ToAlert()
	created, err := s.store.CreateAlert(ctx, alert)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return created, nil
}

// GetAlert 获取告警
func (s *AlertService) GetAlert(ctx context.Context, id int64) (*model.Alert, error) {
	alert, err := s.store.GetAlert(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return nil, ErrAlertNotFound
	}
	return alert, nil
}

// ListAlerts 列出所有告警
func (s *AlertService) ListAlerts(ctx context.Context) ([]*model.Alert, error) {
	alerts, err := s.store.ListAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	sort.Slice(alerts, func(i, j int) bool {
		return alertLevelPriority(alerts[i].Level) > alertLevelPriority(alerts[j].Level)
	})

	return alerts, nil
}

// ListActiveAlerts 列出活跃告警
func (s *AlertService) ListActiveAlerts(ctx context.Context) ([]*model.Alert, error) {
	alerts, err := s.store.ListAlertsByStatus(ctx, model.AlertStatusActive)
	if err != nil {
		return nil, fmt.Errorf("failed to list active alerts: %w", err)
	}

	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].CreatedAt.After(alerts[j].CreatedAt)
	})

	return alerts, nil
}

// ListAlertsByDam 按坝体列出告警
func (s *AlertService) ListAlertsByDam(ctx context.Context, damID int64) ([]*model.Alert, error) {
	alerts, err := s.store.ListAlertsByDam(ctx, damID)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by dam: %v", err)
	}
	return alerts, nil
}

// ResolveAlert 解决告警
func (s *AlertService) ResolveAlert(ctx context.Context, id int64, resolvedBy string) error {
	alert, err := s.store.GetAlert(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return ErrAlertNotFound
	}

	if alert.Status == model.AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	err = s.store.ResolveAlert(ctx, id, resolvedBy)
	if err != nil {
		return fmt.Errorf("failed to resolve alert: %w", err)
	}

	return nil
}

// AcknowledgeAlert 确认告警
func (s *AlertService) AcknowledgeAlert(ctx context.Context, id int64, ackBy string) error {
	alert, err := s.store.GetAlert(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}
	if alert == nil {
		return ErrAlertNotFound
	}

	err = s.store.AcknowledgeAlert(ctx, id, ackBy)
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

// CheckAndCreateAlert 检查阈值并创建告警
func (s *AlertService) CheckAndCreateAlert(ctx context.Context, damID, pointID int64, readingType string, value, threshold float64) (*model.Alert, error) {
	if value <= threshold {
		return nil, nil
	}

	level := model.AlertLevelWarning
	if value > threshold*1.5 {
		level = model.AlertLevelDanger
	}
	if value > threshold*2 {
		level = model.AlertLevelCritical
	}

	alert := &model.Alert{
		DamID:        damID,
		PointID:      pointID,
		Level:        level,
		Status:       model.AlertStatusActive,
		Title:        fmt.Sprintf("%s 超限告警", readingType),
		Message:      fmt.Sprintf("%s 当前值 %.2f 超过阈值 %.2f", readingType, value, threshold),
		Threshold:    threshold,
		CurrentValue: value,
		ReadingType:  readingType,
	}

	created, err := s.store.CreateAlert(ctx, alert)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return created, nil
}

// CheckSeepageReading 检查渗流读数并生成告警
func (s *AlertService) CheckSeepageReading(ctx context.Context, damID, pointID int64, reading *model.SeepageReading) ([]*model.Alert, error) {
	var alerts []*model.Alert

	if alertInput, shouldAlert := model.ShouldAlertSeepage(reading, nil); shouldAlert {
		alertInput.DamID = damID
		alertInput.PointID = pointID
		alert := alertInput.ToAlert()
		created, err := s.store.CreateAlert(ctx, alert)
		if err != nil {
			return nil, fmt.Errorf("failed to create seepage alert: %v", err)
		}
		alerts = append(alerts, created)
	}

	return alerts, nil
}

// CheckDisplacementReading 检查位移读数并生成告警
func (s *AlertService) CheckDisplacementReading(ctx context.Context, damID, pointID int64, reading *model.DisplacementReading) ([]*model.Alert, error) {
	var alerts []*model.Alert

	if alertInput, shouldAlert := model.ShouldAlertDisplacement(reading, nil); shouldAlert {
		alertInput.DamID = damID
		alertInput.PointID = pointID
		alert := alertInput.ToAlert()
		created, err := s.store.CreateAlert(ctx, alert)
		if err != nil {
			return nil, fmt.Errorf("failed to create displacement alert: %v", err)
		}
		alerts = append(alerts, created)
	}

	return alerts, nil
}

// CheckPorePressureReading 检查孔隙水压力读数并生成告警
func (s *AlertService) CheckPorePressureReading(ctx context.Context, damID, pointID int64, reading *model.PorePressureReading) ([]*model.Alert, error) {
	var alerts []*model.Alert

	if alertInput, shouldAlert := model.ShouldAlertPorePressure(reading, nil); shouldAlert {
		alertInput.DamID = damID
		alertInput.PointID = pointID
		alert := alertInput.ToAlert()
		created, err := s.store.CreateAlert(ctx, alert)
		if err != nil {
			return nil, fmt.Errorf("failed to create pore pressure alert: %v", err)
		}
		alerts = append(alerts, created)
	}

	return alerts, nil
}

// GetAlertStatistics 获取告警统计
func (s *AlertService) GetAlertStatistics(ctx context.Context) (*model.AlertStatistics, error) {
	alerts, err := s.store.ListAlerts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts for statistics: %v", err)
	}

	stats := &model.AlertStatistics{}
	for _, a := range alerts {
		stats.Total++
		switch a.Status {
		case model.AlertStatusActive:
			stats.Active++
		case model.AlertStatusResolved:
			stats.Resolved++
		}
		switch a.Level {
		case model.AlertLevelCritical:
			stats.Critical++
		case model.AlertLevelDanger:
			stats.Danger++
		case model.AlertLevelWarning:
			stats.Warning++
		case model.AlertLevelInfo:
			stats.Info++
		}
	}

	return stats, nil
}

// alertLevelPriority 获取告警等级优先级
func alertLevelPriority(level model.AlertLevel) int {
	switch level {
	case model.AlertLevelCritical:
		return 4
	case model.AlertLevelDanger:
		return 3
	case model.AlertLevelWarning:
		return 2
	case model.AlertLevelInfo:
		return 1
	default:
		return 0
	}
}
