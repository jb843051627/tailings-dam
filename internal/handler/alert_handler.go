package handler

import (
	"errors"
	"net/http"

	"tailings-dam/internal/model"
	"tailings-dam/internal/service"
)

// HandleListAlerts 列出所有告警
func (h *Handler) HandleListAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	if status != "" {
		alerts, err := h.alertService.ListActiveAlerts(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, alerts)
		return
	}

	alerts, err := h.alertService.ListAlerts(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, alerts)
}

// HandleCreateAlert 创建告警
func (h *Handler) HandleCreateAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.AlertInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	alert, err := h.alertService.CreateAlert(ctx, &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, alert)
}

// HandleGetAlert 获取告警详情
func (h *Handler) HandleGetAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	alert, err := h.alertService.GetAlert(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			respondError(w, http.StatusNotFound, "alert not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, alert)
}

// HandleResolveAlert 解决告警
func (h *Handler) HandleResolveAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.AlertResolveInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.alertService.ResolveAlert(ctx, id, input.ResolvedBy)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			respondError(w, http.StatusNotFound, "alert not found")
			return
		}
		if errors.Is(err, service.ErrAlertAlreadyResolved) {
			respondError(w, http.StatusBadRequest, "alert already resolved")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "alert resolved"})
}

// HandleAcknowledgeAlert 确认告警
func (h *Handler) HandleAcknowledgeAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.AlertAckInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.alertService.AcknowledgeAlert(ctx, id, input.AcknowledgedBy)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			respondError(w, http.StatusNotFound, "alert not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "alert acknowledged"})
}
