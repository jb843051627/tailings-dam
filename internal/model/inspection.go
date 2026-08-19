package model

import "time"

// InspectionStatus 巡检状态
type InspectionStatus string

const (
	InspectionStatusPending   InspectionStatus = "pending"
	InspectionStatusInProgress InspectionStatus = "in_progress"
	InspectionStatusCompleted InspectionStatus = "completed"
	InspectionStatusOverdue   InspectionStatus = "overdue"
	InspectionStatusCancelled InspectionStatus = "cancelled"
)

// InspectionPriority 巡检优先级
type InspectionPriority string

const (
	InspectionPriorityLow    InspectionPriority = "low"
	InspectionPriorityNormal InspectionPriority = "normal"
	InspectionPriorityHigh   InspectionPriority = "high"
	InspectionPriorityUrgent InspectionPriority = "urgent"
)

// Inspection 巡检实体
type Inspection struct {
	ID            int64             `json:"id"`
	DamID         int64             `json:"dam_id"`
	Inspector     string            `json:"inspector"`
	Title         string            `json:"title"`
	ScheduledDate time.Time         `json:"scheduled_date"`
	CompletedDate time.Time         `json:"completed_date"`
	Findings      string            `json:"findings"`
	Status        InspectionStatus `json:"status"`
	Priority      InspectionPriority `json:"priority"`
	Items         []InspectionItem  `json:"items"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// InspectionItem 巡检检查项
type InspectionItem struct {
	ID            int64  `json:"id"`
	InspectionID  int64  `json:"inspection_id"`
	Name          string `json:"name"`
	Result        string `json:"result"`
	Remark        string `json:"remark"`
	Pass          bool   `json:"pass"`
}

// InspectionInput 创建/更新巡检的输入
type InspectionInput struct {
	DamID         int64               `json:"dam_id"`
	Inspector     string              `json:"inspector"`
	Title         string              `json:"title"`
	ScheduledDate string              `json:"scheduled_date"`
	Priority      InspectionPriority  `json:"priority"`
	Findings      string              `json:"findings"`
}

// ToInspection 将输入转换为 Inspection 实体
func (input *InspectionInput) ToInspection() *Inspection {
	var scheduledDate time.Time
	if input.ScheduledDate != "" {
		scheduledDate, _ = time.Parse(time.RFC3339, input.ScheduledDate)
	}
	priority := input.Priority
	if priority == "" {
		priority = InspectionPriorityNormal
	}
	return &Inspection{
		DamID:         input.DamID,
		Inspector:     input.Inspector,
		Title:         input.Title,
		ScheduledDate: scheduledDate,
		Status:        InspectionStatusPending,
		Priority:      priority,
		Findings:      input.Findings,
	}
}

// InspectionWithDam 带坝体名称的巡检
type InspectionWithDam struct {
	Inspection
	DamName string `json:"dam_name"`
}

// InspectionStatus 巡检状态中文标签
func (s InspectionStatus) Label() string {
	switch s {
	case InspectionStatusPending:
		return "待巡检"
	case InspectionStatusInProgress:
		return "巡检中"
	case InspectionStatusCompleted:
		return "已完成"
	case InspectionStatusOverdue:
		return "已逾期"
	case InspectionStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// InspectionPriority 巡检优先级中文标签
func (p InspectionPriority) Label() string {
	switch p {
	case InspectionPriorityLow:
		return "低"
	case InspectionPriorityNormal:
		return "普通"
	case InspectionPriorityHigh:
		return "高"
	case InspectionPriorityUrgent:
		return "紧急"
	default:
		return "未知"
	}
}

// Validate 验证巡检输入
func (input *InspectionInput) Validate() error {
	if input.DamID <= 0 {
		return &ValidationError{Field: "dam_id", Msg: "坝体ID不能为空"}
	}
	if input.Inspector == "" {
		return &ValidationError{Field: "inspector", Msg: "巡检人不能为空"}
	}
	if input.Title == "" {
		return &ValidationError{Field: "title", Msg: "巡检标题不能为空"}
	}
	if input.ScheduledDate == "" {
		return &ValidationError{Field: "scheduled_date", Msg: "计划日期不能为空"}
	}
	switch input.Priority {
	case InspectionPriorityLow, InspectionPriorityNormal, InspectionPriorityHigh, InspectionPriorityUrgent, "":
	default:
		return &ValidationError{Field: "priority", Msg: "优先级无效"}
	}
	return nil
}

// InspectionCompleteInput 完成巡检输入
type InspectionCompleteInput struct {
	Findings string `json:"findings"`
	Items    []InspectionItem `json:"items"`
}

// InspectionSummary 巡检摘要
type InspectionSummary struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	Completed  int `json:"completed"`
	Overdue    int `json:"overdue"`
	InProgress int `json:"in_progress"`
}

// InspectionQuery 巡检查询参数
type InspectionQuery struct {
	DamID    int64             `json:"dam_id"`
	Status   InspectionStatus  `json:"status"`
	Priority InspectionPriority `json:"priority"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// IsOverdue 检查巡检是否逾期
func (i *Inspection) IsOverdue() bool {
	if i.Status == InspectionStatusCompleted || i.Status == InspectionStatusCancelled {
		return false
	}
	return time.Now().After(i.ScheduledDate)
}

// MarkCompleted 标记巡检为已完成
func (i *Inspection) MarkCompleted(findings string) {
	i.Status = InspectionStatusCompleted
	i.Findings = findings
	i.CompletedDate = time.Now()
	i.UpdatedAt = time.Now()
}

// MarkInProgress 标记巡检为进行中
func (i *Inspection) MarkInProgress() {
	i.Status = InspectionStatusInProgress
	i.UpdatedAt = time.Now()
}

// MarkOverdue 标记巡检为逾期
func (i *Inspection) MarkOverdue() {
	if i.Status == InspectionStatusPending && i.IsOverdue() {
		i.Status = InspectionStatusOverdue
		i.UpdatedAt = time.Now()
	}
}
