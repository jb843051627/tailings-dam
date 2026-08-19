package model

import "time"

// DamStatus 表示尾矿坝的运行状态
type DamStatus string

const (
	DamStatusActive    DamStatus = "active"
	DamStatusInactive  DamStatus = "inactive"
	DamStatusClosed    DamStatus = "closed"
	DamStatusOverflow  DamStatus = "overflow_risk"
	DamStatusEmergency DamStatus = "emergency"
)

// HazardLevel 表示尾矿坝的危险等级
type HazardLevel string

const (
	HazardLevelLow      HazardLevel = "low"
	HazardLevelMedium   HazardLevel = "medium"
	HazardLevelHigh     HazardLevel = "high"
	HazardLevelCritical HazardLevel = "critical"
)

// Dam 尾矿坝实体
type Dam struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Location     string      `json:"location"`
	Province     string      `json:"province"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	Capacity     float64     `json:"capacity"`
	CurrentLevel float64     `json:"current_level"`
	HazardLevel  HazardLevel `json:"hazard_level"`
	Status       DamStatus   `json:"status"`
	Description  string      `json:"description"`
	Operator     string      `json:"operator"`
	ConstructedAt time.Time  `json:"constructed_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// DamInput 创建/更新尾矿坝的输入
type DamInput struct {
	Name          string      `json:"name"`
	Location      string      `json:"location"`
	Province      string      `json:"province"`
	Latitude      float64     `json:"latitude"`
	Longitude     float64     `json:"longitude"`
	Capacity      float64     `json:"capacity"`
	CurrentLevel  float64     `json:"current_level"`
	HazardLevel   HazardLevel `json:"hazard_level"`
	Status        DamStatus   `json:"status"`
	Description   string      `json:"description"`
	Operator      string      `json:"operator"`
	ConstructedAt string      `json:"constructed_at"`
}

// ToDam 将输入转换为 Dam 实体
func (input *DamInput) ToDam() *Dam {
	var constructedAt time.Time
	if input.ConstructedAt != "" {
		constructedAt, _ = time.Parse(time.RFC3339, input.ConstructedAt)
	}
	return &Dam{
		Name:          input.Name,
		Location:      input.Location,
		Province:      input.Province,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		Capacity:      input.Capacity,
		CurrentLevel:  input.CurrentLevel,
		HazardLevel:   input.HazardLevel,
		Status:        input.Status,
		Description:   input.Description,
		Operator:      input.Operator,
		ConstructedAt: constructedAt,
	}
}

// DamSummary 尾矿坝摘要信息
type DamSummary struct {
	ID              int64       `json:"id"`
	Name            string      `json:"name"`
	Location        string      `json:"location"`
	HazardLevel     HazardLevel `json:"hazard_level"`
	Status          DamStatus   `json:"status"`
	AlertCount      int         `json:"alert_count"`
	PointCount      int         `json:"point_count"`
	CapacityUsage   float64     `json:"capacity_usage"`
}

// Validate 验证尾矿坝输入
func (input *DamInput) Validate() error {
	if input.Name == "" {
		return ErrDamNameRequired
	}
	if input.Location == "" {
		return ErrDamLocationRequired
	}
	if input.Capacity <= 0 {
		return ErrDamCapacityInvalid
	}
	switch input.HazardLevel {
	case HazardLevelLow, HazardLevelMedium, HazardLevelHigh, HazardLevelCritical:
	default:
		return ErrDamHazardLevelInvalid
	}
	switch input.Status {
	case DamStatusActive, DamStatusInactive, DamStatusClosed, DamStatusOverflow, DamStatusEmergency:
	default:
		input.Status = DamStatusActive
	}
	return nil
}

// 验证错误
var (
	ErrDamNameRequired       = &ValidationError{Field: "name", Msg: "坝体名称不能为空"}
	ErrDamLocationRequired   = &ValidationError{Field: "location", Msg: "坝体位置不能为空"}
	ErrDamCapacityInvalid    = &ValidationError{Field: "capacity", Msg: "坝体容量必须大于0"}
	ErrDamHazardLevelInvalid = &ValidationError{Field: "hazard_level", Msg: "危险等级无效"}
)

// ValidationError 字段验证错误
type ValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Msg
}
