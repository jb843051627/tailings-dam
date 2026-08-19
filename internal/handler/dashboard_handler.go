package handler

import "net/http"

// DashboardResponse 仪表盘响应
type DashboardResponse struct {
	DamCount          int64      `json:"dam_count"`
	PointCount        int64      `json:"point_count"`
	ActiveAlertCount  int        `json:"active_alert_count"`
	PendingInspection int64      `json:"pending_inspection"`
	AbnormalDrainage  int64      `json:"abnormal_drainage"`
}

// HandleDashboard 仪表盘汇总
func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	damCount, _ := h.damService.GetDamCount(ctx)
	pointCount, _ := h.monitoringService.GetMonitoringPointCount(ctx)
	activeAlerts, _ := h.alertService.ListActiveAlerts(ctx)
	inspSummary, _ := h.inspectionService.GetInspectionSummary(ctx)
	drainageSummary, _ := h.drainageService.GetDrainageSummary(ctx)
	alertStats, _ := h.alertService.GetAlertStatistics(ctx)
	abnormalDrainage, _ := h.drainageService.GetAbnormalDrainageCount(ctx)

	resp := map[string]interface{}{
		"dam_count":          damCount,
		"point_count":        pointCount,
		"active_alert_count": len(activeAlerts),
		"alert_stats":        alertStats,
		"inspection_summary": inspSummary,
		"drainage_summary":   drainageSummary,
		"abnormal_drainage":  abnormalDrainage,
	}

	respondJSON(w, http.StatusOK, resp)
}
