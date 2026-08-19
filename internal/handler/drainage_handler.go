package handler

import (
	"errors"
	"net/http"

	"tailings-dam/internal/model"
	"tailings-dam/internal/service"
)

// HandleListDrainageSystems 列出所有排水系统
func (h *Handler) HandleListDrainageSystems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	damIDStr := r.URL.Query().Get("dam_id")
	if damIDStr != "" {
		damID, err := parseQueryInt64(r, "dam_id")
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid dam_id")
			return
		}
		systems, err := h.drainageService.ListDrainageSystemsByDam(ctx, damID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, systems)
		return
	}

	systems, err := h.drainageService.ListDrainageSystems(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, systems)
}

// HandleCreateDrainageSystem 创建排水系统
func (h *Handler) HandleCreateDrainageSystem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.DrainageInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	system, err := h.drainageService.CreateDrainageSystem(ctx, &input)
	if err != nil {
		if errors.Is(err, service.ErrPointDamNotFound) {
			respondError(w, http.StatusBadRequest, "referenced dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, system)
}

// HandleGetDrainageSystem 获取排水系统详情
func (h *Handler) HandleGetDrainageSystem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	system, err := h.drainageService.GetDrainageSystem(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrDrainageNotFound) {
			respondError(w, http.StatusNotFound, "drainage system not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, system)
}

// HandleUpdateDrainageSystem 更新排水系统
func (h *Handler) HandleUpdateDrainageSystem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.DrainageInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	system, err := h.drainageService.UpdateDrainageSystem(ctx, id, &input)
	if err != nil {
		if errors.Is(err, service.ErrDrainageNotFound) {
			respondError(w, http.StatusNotFound, "drainage system not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, system)
}
