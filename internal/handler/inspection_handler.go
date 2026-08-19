package handler

import (
	"errors"
	"net/http"

	"tailings-dam/internal/model"
	"tailings-dam/internal/service"
)

// HandleListInspections 列出所有巡检
func (h *Handler) HandleListInspections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	damIDStr := r.URL.Query().Get("dam_id")
	if damIDStr != "" {
		damID, err := parseQueryInt64(r, "dam_id")
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid dam_id")
			return
		}
		inspections, err := h.inspectionService.ListInspectionsByDam(ctx, damID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, inspections)
		return
	}

	inspections, err := h.inspectionService.ListInspections(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, inspections)
}

// HandleCreateInspection 创建巡检
func (h *Handler) HandleCreateInspection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.InspectionInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	insp, err := h.inspectionService.CreateInspection(ctx, &input)
	if err != nil {
		if errors.Is(err, service.ErrPointDamNotFound) {
			respondError(w, http.StatusBadRequest, "referenced dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, insp)
}

// HandleGetInspection 获取巡检详情
func (h *Handler) HandleGetInspection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	insp, err := h.inspectionService.GetInspection(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrInspectionNotFound) {
			respondError(w, http.StatusNotFound, "inspection not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, insp)
}

// HandleCompleteInspection 完成巡检
func (h *Handler) HandleCompleteInspection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.InspectionCompleteInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.inspectionService.CompleteInspection(ctx, id, input.Findings)
	if err != nil {
		if errors.Is(err, service.ErrInspectionNotFound) {
			respondError(w, http.StatusNotFound, "inspection not found")
			return
		}
		if errors.Is(err, service.ErrInspectionCompleted) {
			respondError(w, http.StatusBadRequest, "inspection already completed")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "inspection completed"})
}
