package model

import "time"

// SeepageReading 渗流读数
type SeepageReading struct {
	ID         int64     `json:"id"`
	PointID    int64     `json:"point_id"`
	FlowRate   float64   `json:"flow_rate"`
	Turbidity  float64   `json:"turbidity"`
	PH         float64   `json:"ph"`
	Temperature float64  `json:"temperature"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// SeepageReadingInput 渗流读数输入
type SeepageReadingInput struct {
	PointID     int64   `json:"point_id"`
	FlowRate    float64 `json:"flow_rate"`
	Turbidity   float64 `json:"turbidity"`
	PH          float64 `json:"ph"`
	Temperature float64 `json:"temperature"`
	RecordedAt  string  `json:"recorded_at"`
}

// ToSeepageReading 将输入转换为 SeepageReading
func (input *SeepageReadingInput) ToSeepageReading() *SeepageReading {
	var recordedAt time.Time
	if input.RecordedAt != "" {
		recordedAt, _ = time.Parse(time.RFC3339, input.RecordedAt)
	}
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}
	return &SeepageReading{
		PointID:     input.PointID,
		FlowRate:    input.FlowRate,
		Turbidity:   input.Turbidity,
		PH:          input.PH,
		Temperature: input.Temperature,
		RecordedAt:  recordedAt,
	}
}

// DisplacementReading 位移读数
type DisplacementReading struct {
	ID                     int64     `json:"id"`
	PointID                int64     `json:"point_id"`
	HorizontalDisplacement float64   `json:"horizontal_displacement"`
	VerticalDisplacement   float64   `json:"vertical_displacement"`
	CumulativeHorizontal   float64   `json:"cumulative_horizontal"`
	CumulativeVertical     float64   `json:"cumulative_vertical"`
	RecordedAt             time.Time `json:"recorded_at"`
	CreatedAt              time.Time `json:"created_at"`
}

// DisplacementReadingInput 位移读数输入
type DisplacementReadingInput struct {
	PointID                int64   `json:"point_id"`
	HorizontalDisplacement float64 `json:"horizontal_displacement"`
	VerticalDisplacement   float64 `json:"vertical_displacement"`
	CumulativeHorizontal   float64 `json:"cumulative_horizontal"`
	CumulativeVertical     float64 `json:"cumulative_vertical"`
	RecordedAt             string  `json:"recorded_at"`
}

// ToDisplacementReading 将输入转换为 DisplacementReading
func (input *DisplacementReadingInput) ToDisplacementReading() *DisplacementReading {
	var recordedAt time.Time
	if input.RecordedAt != "" {
		recordedAt, _ = time.Parse(time.RFC3339, input.RecordedAt)
	}
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}
	return &DisplacementReading{
		PointID:                input.PointID,
		HorizontalDisplacement: input.HorizontalDisplacement,
		VerticalDisplacement:   input.VerticalDisplacement,
		CumulativeHorizontal:   input.CumulativeHorizontal,
		CumulativeVertical:     input.CumulativeVertical,
		RecordedAt:             recordedAt,
	}
}

// PorePressureReading 孔隙水压力读数
type PorePressureReading struct {
	ID         int64     `json:"id"`
	PointID    int64     `json:"point_id"`
	Pressure   float64   `json:"pressure"`
	Depth       float64   `json:"depth"`
	WaterLevel  float64   `json:"water_level"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// PorePressureReadingInput 孔隙水压力读数输入
type PorePressureReadingInput struct {
	PointID    int64   `json:"point_id"`
	Pressure   float64 `json:"pressure"`
	Depth      float64 `json:"depth"`
	WaterLevel float64 `json:"water_level"`
	RecordedAt string  `json:"recorded_at"`
}

// ToPorePressureReading 将输入转换为 PorePressureReading
func (input *PorePressureReadingInput) ToPorePressureReading() *PorePressureReading {
	var recordedAt time.Time
	if input.RecordedAt != "" {
		recordedAt, _ = time.Parse(time.RFC3339, input.RecordedAt)
	}
	if recordedAt.IsZero() {
		recordedAt = time.Now()
	}
	return &PorePressureReading{
		PointID:    input.PointID,
		Pressure:   input.Pressure,
		Depth:      input.Depth,
		WaterLevel: input.WaterLevel,
		RecordedAt: recordedAt,
	}
}

// BatchReadingInput 批量读数输入
type BatchReadingInput struct {
	SeepageReadings      []SeepageReadingInput      `json:"seepage_readings"`
	DisplacementReadings []DisplacementReadingInput `json:"displacement_readings"`
	PorePressureReadings []PorePressureReadingInput  `json:"pore_pressure_readings"`
}

// ReadingSummary 读数摘要
type ReadingSummary struct {
	PointID    int64   `json:"point_id"`
	PointName  string  `json:"point_name"`
	PointType  string  `json:"point_type"`
	LatestValue float64 `json:"latest_value"`
	LatestTime  string  `json:"latest_time"`
	ReadingsCount int  `json:"readings_count"`
}

// SeepageThreshold 渗流阈值
type SeepageThreshold struct {
	MaxFlowRate  float64
	MaxTurbidity float64
	MinPH        float64
	MaxPH        float64
}

// DisplacementThreshold 位移阈值
type DisplacementThreshold struct {
	MaxHorizontal float64
	MaxVertical   float64
}

// PorePressureThreshold 孔隙水压力阈值
type PorePressureThreshold struct {
	MaxPressure float64
	MaxWaterLevel float64
}

// DefaultSeepageThreshold 默认渗流阈值
var DefaultSeepageThreshold = SeepageThreshold{
	MaxFlowRate:  100.0,
	MaxTurbidity: 50.0,
	MinPH:        6.0,
	MaxPH:        9.0,
}

// DefaultDisplacementThreshold 默认位移阈值
var DefaultDisplacementThreshold = DisplacementThreshold{
	MaxHorizontal: 50.0,
	MaxVertical:   30.0,
}

// DefaultPorePressureThreshold 默认孔隙水压力阈值
var DefaultPorePressureThreshold = PorePressureThreshold{
	MaxPressure:   500.0,
	MaxWaterLevel: 100.0,
}

// ReadingType 读数类型
type ReadingType string

const (
	ReadingTypeSeepage       ReadingType = "seepage"
	ReadingTypeDisplacement  ReadingType = "displacement"
	ReadingTypePorePressure  ReadingType = "pore_pressure"
)

// ReadingQuery 读数查询参数
type ReadingQuery struct {
	PointID    int64     `json:"point_id"`
	Type       ReadingType `json:"type"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

// Validate 验证渗流读数输入
func (input *SeepageReadingInput) Validate() error {
	if input.PointID <= 0 {
		return &ValidationError{Field: "point_id", Msg: "监测点ID不能为空"}
	}
	if input.FlowRate < 0 {
		return &ValidationError{Field: "flow_rate", Msg: "流量不能为负数"}
	}
	if input.Turbidity < 0 {
		return &ValidationError{Field: "turbidity", Msg: "浊度不能为负数"}
	}
	return nil
}

// Validate 验证位移读数输入
func (input *DisplacementReadingInput) Validate() error {
	if input.PointID <= 0 {
		return &ValidationError{Field: "point_id", Msg: "监测点ID不能为空"}
	}
	return nil
}

// Validate 验证孔隙水压力读数输入
func (input *PorePressureReadingInput) Validate() error {
	if input.PointID <= 0 {
		return &ValidationError{Field: "point_id", Msg: "监测点ID不能为空"}
	}
	if input.Pressure < 0 {
		return &ValidationError{Field: "pressure", Msg: "压力不能为负数"}
	}
	return nil
}
