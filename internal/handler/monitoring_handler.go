package handler

import (
	"errors"
	"net/http"

	"tailings-dam/internal/model"
	"tailings-dam/internal/service"
)

// HandleListMonitoringPoints 列出所有监测点
func (h *Handler) HandleListMonitoringPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	damIDStr := r.URL.Query().Get("dam_id")
	if damIDStr != "" {
		damID, err := parseQueryInt64(r, "dam_id")
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid dam_id")
			return
		}
		points, err := h.monitoringService.ListMonitoringPointsByDam(ctx, damID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, points)
		return
	}

	points, err := h.monitoringService.ListMonitoringPoints(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, points)
}

// HandleCreateMonitoringPoint 创建监测点
func (h *Handler) HandleCreateMonitoringPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.MonitoringPointInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	point, err := h.monitoringService.CreateMonitoringPoint(ctx, &input)
	if err != nil {
		if errors.Is(err, service.ErrPointDamNotFound) {
			respondError(w, http.StatusBadRequest, "referenced dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, point)
}

// HandleGetMonitoringPoint 获取监测点详情
func (h *Handler) HandleGetMonitoringPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	point, err := h.monitoringService.GetMonitoringPoint(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrPointNotFound) {
			respondError(w, http.StatusNotFound, "monitoring point not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, point)
}

// HandleUpdateMonitoringPoint 更新监测点
func (h *Handler) HandleUpdateMonitoringPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.MonitoringPointInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	point, err := h.monitoringService.UpdateMonitoringPoint(ctx, id, &input)
	if err != nil {
		if errors.Is(err, service.ErrPointNotFound) {
			respondError(w, http.StatusNotFound, "monitoring point not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, point)
}

// HandleDeleteMonitoringPoint 删除监测点
func (h *Handler) HandleDeleteMonitoringPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.monitoringService.DeleteMonitoringPoint(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrPointNotFound) {
			respondError(w, http.StatusNotFound, "monitoring point not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "monitoring point deleted"})
}
