package model

import "time"

// PointType 监测点类型
type PointType string

const (
	PointTypeSeepage       PointType = "seepage"
	PointTypeDisplacement  PointType = "displacement"
	PointTypePorePressure  PointType = "pore_pressure"
	PointTypePhreaticSurface PointType = "phreatic_surface"
)

// PointStatus 监测点状态
type PointStatus string

const (
	PointStatusActive   PointStatus = "active"
	PointStatusInactive PointStatus = "inactive"
	PointStatusFaulty   PointStatus = "faulty"
	PointStatusMaintain PointStatus = "maintenance"
)

// MonitoringPoint 监测点实体
type MonitoringPoint struct {
	ID          int64       `json:"id"`
	DamID       int64       `json:"dam_id"`
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Type        PointType   `json:"type"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	Elevation   float64     `json:"elevation"`
	Status      PointStatus `json:"status"`
	Description string      `json:"description"`
	LastReading time.Time   `json:"last_reading"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// MonitoringPointInput 创建/更新监测点的输入
type MonitoringPointInput struct {
	DamID       int64     `json:"dam_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Type        PointType `json:"type"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Elevation   float64   `json:"elevation"`
	Status      PointStatus `json:"status"`
	Description string    `json:"description"`
}

// ToMonitoringPoint 将输入转换为 MonitoringPoint 实体
func (input *MonitoringPointInput) ToMonitoringPoint() *MonitoringPoint {
	status := input.Status
	if status == "" {
		status = PointStatusActive
	}
	return &MonitoringPoint{
		DamID:       input.DamID,
		Name:        input.Name,
		Code:        input.Code,
		Type:        input.Type,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		Elevation:   input.Elevation,
		Status:      status,
		Description: input.Description,
	}
}

// MonitoringPointWithType 带坝体名称的监测点
type MonitoringPointWithType struct {
	MonitoringPoint
	DamName string `json:"dam_name"`
}

// Validate 验证监测点输入
func (input *MonitoringPointInput) Validate() error {
	if input.DamID <= 0 {
		return ErrPointDamIDRequired
	}
	if input.Name == "" {
		return ErrPointNameRequired
	}
	if input.Code == "" {
		return ErrPointCodeRequired
	}
	switch input.Type {
	case PointTypeSeepage, PointTypeDisplacement, PointTypePorePressure, PointTypePhreaticSurface:
	default:
		return ErrPointTypeInvalid
	}
	switch input.Status {
	case PointStatusActive, PointStatusInactive, PointStatusFaulty, PointStatusMaintain, "":
	default:
		return ErrPointStatusInvalid
	}
	return nil
}

// 监测点验证错误
var (
	ErrPointDamIDRequired  = &ValidationError{Field: "dam_id", Msg: "坝体ID不能为空"}
	ErrPointNameRequired   = &ValidationError{Field: "name", Msg: "监测点名称不能为空"}
	ErrPointCodeRequired   = &ValidationError{Field: "code", Msg: "监测点编号不能为空"}
	ErrPointTypeInvalid    = &ValidationError{Field: "type", Msg: "监测点类型无效"}
	ErrPointStatusInvalid  = &ValidationError{Field: "status", Msg: "监测点状态无效"}
)

// PointTypeLabel 获取监测点类型的中文标签
func (t PointType) Label() string {
	switch t {
	case PointTypeSeepage:
		return "渗流监测"
	case PointTypeDisplacement:
		return "位移监测"
	case PointTypePorePressure:
		return "孔隙水压力监测"
	case PointTypePhreaticSurface:
		return "浸润线监测"
	default:
		return "未知类型"
	}
}

// PointStatusLabel 获取监测点状态的中文标签
func (s PointStatus) Label() string {
	switch s {
	case PointStatusActive:
		return "正常"
	case PointStatusInactive:
		return "停用"
	case PointStatusFaulty:
		return "故障"
	case PointStatusMaintain:
		return "维护中"
	default:
		return "未知状态"
	}
}
