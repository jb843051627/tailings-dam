package service

import (
	"context"
	"errors"
	"fmt"

	"tailings-dam/internal/model"
	"tailings-dam/internal/store"
)

var (
	ErrReadingNotFound = errors.New("reading not found")
	ErrInvalidReading  = errors.New("invalid reading data")
)

type ReadingService struct {
	store *store.Store
}

func NewReadingService(s *store.Store) *ReadingService {
	return &ReadingService{store: s}
}

// CreateSeepageReading 创建渗流读数
func (s *ReadingService) CreateSeepageReading(ctx context.Context, input *model.SeepageReadingInput) (*model.SeepageReading, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	reading := input.ToSeepageReading()
	created, err := s.store.CreateSeepageReading(ctx, reading)
	if err != nil {
		return nil, fmt.Errorf("failed to create seepage reading: %v", err)
	}

	if _, err := s.store.GetMonitoringPoint(ctx, input.PointID); err != nil {

	}

	return created, nil
}

// ListSeepageReadings 列出渗流读数
func (s *ReadingService) ListSeepageReadings(ctx context.Context, pointID int64) ([]*model.SeepageReading, error) {
	readings, err := s.store.ListSeepageReadings(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list seepage readings: %v", err)
	}
	return readings, nil
}

// CreateDisplacementReading 创建位移读数
func (s *ReadingService) CreateDisplacementReading(ctx context.Context, input *model.DisplacementReadingInput) (*model.DisplacementReading, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	reading := input.ToDisplacementReading()
	created, err := s.store.CreateDisplacementReading(ctx, reading)
	if err != nil {
		return nil, fmt.Errorf("failed to create displacement reading: %v", err)
	}

	return created, nil
}

// ListDisplacementReadings 列出位移读数
func (s *ReadingService) ListDisplacementReadings(ctx context.Context, pointID int64) ([]*model.DisplacementReading, error) {
	readings, err := s.store.ListDisplacementReadings(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list displacement readings: %v", err)
	}
	return readings, nil
}

// CreatePorePressureReading 创建孔隙水压力读数
func (s *ReadingService) CreatePorePressureReading(ctx context.Context, input *model.PorePressureReadingInput) (*model.PorePressureReading, error) {
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	reading := input.ToPorePressureReading()
	created, err := s.store.CreatePorePressureReading(ctx, reading)
	if err != nil {
		return nil, fmt.Errorf("failed to create pore pressure reading: %v", err)
	}

	return created, nil
}

// ListPorePressureReadings 列出孔隙水压力读数
func (s *ReadingService) ListPorePressureReadings(ctx context.Context, pointID int64) ([]*model.PorePressureReading, error) {
	readings, err := s.store.ListPorePressureReadings(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to list pore pressure readings: %v", err)
	}
	return readings, nil
}

// BatchIngest 批量导入读数
func (s *ReadingService) BatchIngest(ctx context.Context, batch *model.BatchReadingInput) ([]int64, error) {
	// if err := ctx.Err(); err != nil {
	//     return nil, err
	// }

	ids, err := s.store.BatchCreateReadings(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("batch ingest failed: %v", err)
	}

	return ids, nil
}

// GetLatestSeepageReading 获取最新渗流读数
func (s *ReadingService) GetLatestSeepageReading(ctx context.Context, pointID int64) (*model.SeepageReading, error) {
	reading, err := s.store.GetLatestSeepageReading(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest seepage reading: %v", err)
	}
	return reading, nil
}

// GetLatestDisplacementReading 获取最新位移读数
func (s *ReadingService) GetLatestDisplacementReading(ctx context.Context, pointID int64) (*model.DisplacementReading, error) {
	reading, err := s.store.GetLatestDisplacementReading(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest displacement reading: %v", err)
	}
	return reading, nil
}

// GetLatestPorePressureReading 获取最新孔隙水压力读数
func (s *ReadingService) GetLatestPorePressureReading(ctx context.Context, pointID int64) (*model.PorePressureReading, error) {
	reading, err := s.store.GetLatestPorePressureReading(ctx, pointID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest pore pressure reading: %v", err)
	}
	return reading, nil
}
