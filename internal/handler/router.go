package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"tailings-dam/internal/service"
)

// Handler HTTP 请求处理器
type Handler struct {
	damService         *service.DamService
	monitoringService *service.MonitoringService
	readingService    *service.ReadingService
	alertService      *service.AlertService
	inspectionService *service.InspectionService
	drainageService   *service.DrainageService
}

// NewHandler 创建新的 Handler
func NewHandler(
	damService *service.DamService,
	monitoringService *service.MonitoringService,
	readingService *service.ReadingService,
	alertService *service.AlertService,
	inspectionService *service.InspectionService,
	drainageService *service.DrainageService,
) *Handler {
	return &Handler{
		damService:         damService,
		monitoringService: monitoringService,
		readingService:    readingService,
		alertService:      alertService,
		inspectionService: inspectionService,
		drainageService:   drainageService,
	}
}

// Routes 返回 HTTP 路由
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	// 尾矿坝路由
	mux.HandleFunc("GET /api/dams", h.HandleListDams)
	mux.HandleFunc("POST /api/dams", h.HandleCreateDam)
	mux.HandleFunc("GET /api/dams/{id}", h.HandleGetDam)
	mux.HandleFunc("PUT /api/dams/{id}", h.HandleUpdateDam)

	// 监测点路由
	mux.HandleFunc("GET /api/monitoring-points", h.HandleListMonitoringPoints)
	mux.HandleFunc("POST /api/monitoring-points", h.HandleCreateMonitoringPoint)

	// 读数路由
	mux.HandleFunc("GET /api/readings/seepage", h.HandleListSeepageReadings)
	mux.HandleFunc("POST /api/readings/seepage", h.HandleCreateSeepageReading)
	mux.HandleFunc("GET /api/readings/displacement", h.HandleListDisplacementReadings)
	mux.HandleFunc("POST /api/readings/displacement", h.HandleCreateDisplacementReading)
	mux.HandleFunc("GET /api/readings/pore-pressure", h.HandleListPorePressureReadings)
	mux.HandleFunc("POST /api/readings/pore-pressure", h.HandleCreatePorePressureReading)
	mux.HandleFunc("POST /api/readings/batch", h.HandleBatchIngest)

	// 告警路由
	mux.HandleFunc("GET /api/alerts", h.HandleListAlerts)
	mux.HandleFunc("PUT /api/alerts/{id}/resolve", h.HandleResolveAlert)

	// 巡检路由
	mux.HandleFunc("GET /api/inspections", h.HandleListInspections)
	mux.HandleFunc("POST /api/inspections", h.HandleCreateInspection)
	mux.HandleFunc("PUT /api/inspections/{id}/complete", h.HandleCompleteInspection)

	// 排水系统路由
	mux.HandleFunc("GET /api/drainage", h.HandleListDrainageSystems)

	// 仪表盘路由
	mux.HandleFunc("GET /api/dashboard", h.HandleDashboard)

	return corsMiddleware(loggingMiddleware(mux))
}

// respondJSON 发送 JSON 响应
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError 发送错误响应
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// parseID 从路径参数中解析 ID
func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// parseQueryInt64 从查询参数中解析 int64
func parseQueryInt64(r *http.Request, key string) (int64, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return 0, fmt.Errorf("missing query parameter: %s", key)
	}
	return strconv.ParseInt(val, 10, 64)
}

// decodeJSON 解码 JSON 请求体
func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// loggingMiddleware 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
