package handler

import (
	"net/http"

	"tailings-dam/internal/model"
)

// HandleListSeepageReadings 列出渗流读数
func (h *Handler) HandleListSeepageReadings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pointID, err := parseQueryInt64(r, "point_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid or missing point_id")
		return
	}

	readings, err := h.readingService.ListSeepageReadings(ctx, pointID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, readings)
}

// HandleCreateSeepageReading 创建渗流读数
func (h *Handler) HandleCreateSeepageReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.SeepageReadingInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reading, err := h.readingService.CreateSeepageReading(ctx, &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, reading)
}

// HandleListDisplacementReadings 列出位移读数
func (h *Handler) HandleListDisplacementReadings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pointID, err := parseQueryInt64(r, "point_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid or missing point_id")
		return
	}

	readings, err := h.readingService.ListDisplacementReadings(ctx, pointID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, readings)
}

// HandleCreateDisplacementReading 创建位移读数
func (h *Handler) HandleCreateDisplacementReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.DisplacementReadingInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reading, err := h.readingService.CreateDisplacementReading(ctx, &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, reading)
}

// HandleListPorePressureReadings 列出孔隙水压力读数
func (h *Handler) HandleListPorePressureReadings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pointID, err := parseQueryInt64(r, "point_id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid or missing point_id")
		return
	}

	readings, err := h.readingService.ListPorePressureReadings(ctx, pointID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, readings)
}

// HandleCreatePorePressureReading 创建孔隙水压力读数
func (h *Handler) HandleCreatePorePressureReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.PorePressureReadingInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reading, err := h.readingService.CreatePorePressureReading(ctx, &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, reading)
}

// HandleBatchIngest 批量导入读数
func (h *Handler) HandleBatchIngest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var batch model.BatchReadingInput
	if err := decodeJSON(r, &batch); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ids, err := h.readingService.BatchIngest(ctx, &batch)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"created_count": len(ids),
		"ids":          ids,
	})
}
