package model

import "time"

// DrainageType 排水系统类型
type DrainageType string

const (
	DrainageTypeMain       DrainageType = "main"
	DrainageTypeSecondary  DrainageType = "secondary"
	DrainageTypeEmergency  DrainageType = "emergency"
	DrainageTypeSeepage    DrainageType = "seepage_collection"
	DrainageTypeDischarge  DrainageType = "discharge"
)

// DrainageStatus 排水系统状态
type DrainageStatus string

const (
	DrainageStatusNormal   DrainageStatus = "normal"
	DrainageStatusBlocked  DrainageStatus = "blocked"
	DrainageStatusOverflow DrainageStatus = "overflow"
	DrainageStatusDamaged  DrainageStatus = "damaged"
	DrainageStatusMaintain DrainageStatus = "maintenance"
	DrainageStatusOffline  DrainageStatus = "offline"
)

// DrainageSystem 排水系统实体
type DrainageSystem struct {
	ID             int64          `json:"id"`
	DamID          int64          `json:"dam_id"`
	Name           string         `json:"name"`
	Type           DrainageType   `json:"type"`
	Status         DrainageStatus `json:"status"`
	DesignFlow     float64        `json:"design_flow"`
	ActualFlow     float64        `json:"actual_flow"`
	Diameter       float64        `json:"diameter"`
	Length         float64        `json:"length"`
	Material       string         `json:"material"`
	LastInspection time.Time      `json:"last_inspection"`
	NextInspection time.Time     `json:"next_inspection"`
	Notes          string         `json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// DrainageInput 创建/更新排水系统的输入
type DrainageInput struct {
	DamID          int64          `json:"dam_id"`
	Name           string         `json:"name"`
	Type           DrainageType   `json:"type"`
	Status         DrainageStatus `json:"status"`
	DesignFlow     float64        `json:"design_flow"`
	ActualFlow     float64        `json:"actual_flow"`
	Diameter       float64        `json:"diameter"`
	Length         float64        `json:"length"`
	Material       string         `json:"material"`
	LastInspection string         `json:"last_inspection"`
	NextInspection string         `json:"next_inspection"`
	Notes          string         `json:"notes"`
}

// ToDrainageSystem 将输入转换为 DrainageSystem 实体
func (input *DrainageInput) ToDrainageSystem() *DrainageSystem {
	var prevInspection, nextInspection time.Time
	if input.LastInspection != "" {
		prevInspection, _ = time.Parse(time.RFC3339, input.LastInspection)
	}
	if input.NextInspection != "" {
		nextInspection, _ = time.Parse(time.RFC3339, input.NextInspection)
	}
	status := input.Status
	if status == "" {
		status = DrainageStatusNormal
	}
	return &DrainageSystem{
		DamID:          input.DamID,
		Name:           input.Name,
		Type:           input.Type,
		Status:         status,
		DesignFlow:     input.DesignFlow,
		ActualFlow:     input.ActualFlow,
		Diameter:       input.Diameter,
		Length:         input.Length,
		Material:       input.Material,
		LastInspection: prevInspection,
		NextInspection: nextInspection,
		Notes:          input.Notes,
	}
}

// DrainageWithDam 带坝体名称的排水系统
type DrainageWithDam struct {
	DrainageSystem
	DamName string `json:"dam_name"`
}

// DrainageStatusLabel 排水系统状态中文标签
func (s DrainageStatus) Label() string {
	switch s {
	case DrainageStatusNormal:
		return "正常"
	case DrainageStatusBlocked:
		return "堵塞"
	case DrainageStatusOverflow:
		return "溢流"
	case DrainageStatusDamaged:
		return "损坏"
	case DrainageStatusMaintain:
		return "维护中"
	case DrainageStatusOffline:
		return "离线"
	default:
		return "未知"
	}
}

// DrainageTypeLabel 排水系统类型中文标签
func (t DrainageType) Label() string {
	switch t {
	case DrainageTypeMain:
		return "主排水"
	case DrainageTypeSecondary:
		return "辅助排水"
	case DrainageTypeEmergency:
		return "应急排水"
	case DrainageTypeSeepage:
		return "渗流收集"
	case DrainageTypeDischarge:
		return "排放系统"
	default:
		return "未知类型"
	}
}

// Validate 验证排水系统输入
func (input *DrainageInput) Validate() error {
	if input.DamID <= 0 {
		return &ValidationError{Field: "dam_id", Msg: "坝体ID不能为空"}
	}
	if input.Name == "" {
		return &ValidationError{Field: "name", Msg: "排水系统名称不能为空"}
	}
	switch input.Type {
	case DrainageTypeMain, DrainageTypeSecondary, DrainageTypeEmergency,
		DrainageTypeSeepage, DrainageTypeDischarge:
	default:
		return &ValidationError{Field: "type", Msg: "排水系统类型无效"}
	}
	switch input.Status {
	case DrainageStatusNormal, DrainageStatusBlocked, DrainageStatusOverflow,
		DrainageStatusDamaged, DrainageStatusMaintain, DrainageStatusOffline, "":
	default:
		return &ValidationError{Field: "status", Msg: "排水系统状态无效"}
	}
	if input.DesignFlow < 0 {
		return &ValidationError{Field: "design_flow", Msg: "设计流量不能为负数"}
	}
	return nil
}

// DrainageSummary 排水系统摘要
type DrainageSummary struct {
	Total      int `json:"total"`
	Normal     int `json:"normal"`
	Blocked    int `json:"blocked"`
	Overflow   int `json:"overflow"`
	Damaged    int `json:"damaged"`
	Maintain   int `json:"maintain"`
	Offline    int `json:"offline"`
}

// DrainageQuery 排水系统查询参数
type DrainageQuery struct {
	DamID  int64          `json:"dam_id"`
	Type   DrainageType   `json:"type"`
	Status DrainageStatus `json:"status"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// IsInspectionDue 检查排水系统是否需要巡检
func (d *DrainageSystem) IsInspectionDue() bool {
	if d.NextInspection.IsZero() {
		return false
	}
	return time.Now().After(d.NextInspection)
}

// FlowEfficiency 计算排水效率
func (d *DrainageSystem) FlowEfficiency() float64 {
	if d.DesignFlow <= 0 {
		return 0
	}
	efficiency := (d.ActualFlow / d.DesignFlow) * 100
	if efficiency > 100 {
		efficiency = 100
	}
	return efficiency
}

// IsAbnormal 检查排水系统是否处于异常状态
func (d *DrainageSystem) IsAbnormal() bool {
	switch d.Status {
	case DrainageStatusBlocked, DrainageStatusOverflow,
		DrainageStatusDamaged, DrainageStatusOffline:
		return true
	default:
		return false
	}
}

// MarkInspected 标记排水系统已巡检
func (d *DrainageSystem) MarkInspected(inspectionTime time.Time) {
	d.LastInspection = inspectionTime
	d.UpdatedAt = time.Now()
}
