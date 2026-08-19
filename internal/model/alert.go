package model

import "time"

// AlertLevel 告警等级
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelDanger   AlertLevel = "danger"
	AlertLevelCritical AlertLevel = "critical"
)

// AlertStatus 告警状态
type AlertStatus string

const (
	AlertStatusActive    AlertStatus = "active"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved  AlertStatus = "resolved"
	AlertStatusSuppressed AlertStatus = "suppressed"
)

// Alert 告警实体
type Alert struct {
	ID           int64       `json:"id"`
	DamID        int64       `json:"dam_id"`
	PointID      int64       `json:"point_id"`
	Level        AlertLevel  `json:"level"`
	Status       AlertStatus `json:"status"`
	Title        string      `json:"title"`
	Message      string      `json:"message"`
	Threshold    float64     `json:"threshold"`
	CurrentValue float64     `json:"current_value"`
	ReadingType  string      `json:"reading_type"`
	AcknowledgedBy string    `json:"acknowledged_by"`
	AcknowledgedAt time.Time `json:"acknowledged_at"`
	ResolvedBy   string      `json:"resolved_by"`
	ResolvedAt   time.Time   `json:"resolved_at"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// AlertInput 创建告警的输入
type AlertInput struct {
	DamID        int64      `json:"dam_id"`
	PointID      int64      `json:"point_id"`
	Level        AlertLevel `json:"level"`
	Title        string     `json:"title"`
	Message      string     `json:"message"`
	Threshold    float64    `json:"threshold"`
	CurrentValue float64    `json:"current_value"`
	ReadingType  string     `json:"reading_type"`
}

// ToAlert 将输入转换为 Alert 实体
func (input *AlertInput) ToAlert() *Alert {
	level := input.Level
	if level == "" {
		level = AlertLevelWarning
	}
	return &Alert{
		DamID:        input.DamID,
		PointID:      input.PointID,
		Level:        level,
		Status:       AlertStatusActive,
		Title:        input.Title,
		Message:      input.Message,
		Threshold:    input.Threshold,
		CurrentValue: input.CurrentValue,
		ReadingType:  input.ReadingType,
	}
}

// AlertWithDetails 带详细信息的告警
type AlertWithDetails struct {
	Alert
	DamName   string `json:"dam_name"`
	PointName string `json:"point_name"`
	PointType string `json:"point_type"`
}

// AlertStatistics 告警统计
type AlertStatistics struct {
	Total      int `json:"total"`
	Active     int `json:"active"`
	Resolved   int `json:"resolved"`
	Critical   int `json:"critical"`
	Danger     int `json:"danger"`
	Warning    int `json:"warning"`
	Info       int `json:"info"`
}

// AlertLevelLabel 获取告警等级的中文标签
func (l AlertLevel) Label() string {
	switch l {
	case AlertLevelInfo:
		return "信息"
	case AlertLevelWarning:
		return "警告"
	case AlertLevelDanger:
		return "危险"
	case AlertLevelCritical:
		return "紧急"
	default:
		return "未知"
	}
}

// AlertStatusLabel 获取告警状态的中文标签
func (s AlertStatus) Label() string {
	switch s {
	case AlertStatusActive:
		return "活动中"
	case AlertStatusAcknowledged:
		return "已确认"
	case AlertStatusResolved:
		return "已解决"
	case AlertStatusSuppressed:
		return "已抑制"
	default:
		return "未知"
	}
}

// Validate 验证告警输入
func (input *AlertInput) Validate() error {
	if input.DamID <= 0 {
		return &ValidationError{Field: "dam_id", Msg: "坝体ID不能为空"}
	}
	if input.Message == "" {
		return &ValidationError{Field: "message", Msg: "告警消息不能为空"}
	}
	switch input.Level {
	case AlertLevelInfo, AlertLevelWarning, AlertLevelDanger, AlertLevelCritical, "":
	default:
		return &ValidationError{Field: "level", Msg: "告警等级无效"}
	}
	return nil
}

// AlertAckInput 确认告警输入
type AlertAckInput struct {
	AcknowledgedBy string `json:"acknowledged_by"`
}

// AlertResolveInput 解决告警输入
type AlertResolveInput struct {
	ResolvedBy string `json:"resolved_by"`
	Note      string `json:"note"`
}

// AlertFilter 告警过滤参数
type AlertFilter struct {
	DamID   int64       `json:"dam_id"`
	Level   AlertLevel  `json:"level"`
	Status  AlertStatus `json:"status"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
}

// ShouldAlert 检查渗流读数是否应该触发告警
func ShouldAlertSeepage(reading *SeepageReading, threshold *SeepageThreshold) (*AlertInput, bool) {
	if threshold == nil {
		threshold = &DefaultSeepageThreshold
	}
	if reading.FlowRate > threshold.MaxFlowRate {
		return &AlertInput{
			Level:        AlertLevelDanger,
			Title:        "渗流量超限",
			Message:      "渗流量超过阈值限制",
			Threshold:    threshold.MaxFlowRate,
			CurrentValue: reading.FlowRate,
			ReadingType:  "seepage",
		}, true
	}
	if reading.Turbidity > threshold.MaxTurbidity {
		return &AlertInput{
			Level:        AlertLevelWarning,
			Title:        "浊度超限",
			Message:      "渗流浊度超过阈值限制",
			Threshold:    threshold.MaxTurbidity,
			CurrentValue: reading.Turbidity,
			ReadingType:  "seepage",
		}, true
	}
	if reading.PH < threshold.MinPH || reading.PH > threshold.MaxPH {
		return &AlertInput{
			Level:        AlertLevelWarning,
			Title:        "pH值异常",
			Message:      "渗流pH值超出正常范围",
			Threshold:    threshold.MaxPH,
			CurrentValue: reading.PH,
			ReadingType:  "seepage",
		}, true
	}
	return nil, false
}

// ShouldAlertDisplacement 检查位移读数是否应该触发告警
func ShouldAlertDisplacement(reading *DisplacementReading, threshold *DisplacementThreshold) (*AlertInput, bool) {
	if threshold == nil {
		threshold = &DefaultDisplacementThreshold
	}
	if reading.HorizontalDisplacement > threshold.MaxHorizontal {
		return &AlertInput{
			Level:        AlertLevelDanger,
			Title:        "水平位移超限",
			Message:      "水平位移超过阈值限制",
			Threshold:    threshold.MaxHorizontal,
			CurrentValue: reading.HorizontalDisplacement,
			ReadingType:  "displacement",
		}, true
	}
	if reading.VerticalDisplacement > threshold.MaxVertical {
		return &AlertInput{
			Level:        AlertLevelWarning,
			Title:        "垂直位移超限",
			Message:      "垂直位移超过阈值限制",
			Threshold:    threshold.MaxVertical,
			CurrentValue: reading.VerticalDisplacement,
			ReadingType:  "displacement",
		}, true
	}
	return nil, false
}

// ShouldAlertPorePressure 检查孔隙水压力读数是否应该触发告警
func ShouldAlertPorePressure(reading *PorePressureReading, threshold *PorePressureThreshold) (*AlertInput, bool) {
	if threshold == nil {
		threshold = &DefaultPorePressureThreshold
	}
	if reading.Pressure > threshold.MaxPressure {
		return &AlertInput{
			Level:        AlertLevelCritical,
			Title:        "孔隙水压力超限",
			Message:      "孔隙水压力超过阈值限制",
			Threshold:    threshold.MaxPressure,
			CurrentValue: reading.Pressure,
			ReadingType:  "pore_pressure",
		}, true
	}
	if reading.WaterLevel > threshold.MaxWaterLevel {
		return &AlertInput{
			Level:        AlertLevelDanger,
			Title:        "水位超限",
			Message:      "浸润线水位超过阈值限制",
			Threshold:    threshold.MaxWaterLevel,
			CurrentValue: reading.WaterLevel,
			ReadingType:  "pore_pressure",
		}, true
	}
	return nil, false
}
