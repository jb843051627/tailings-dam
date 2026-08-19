package handler

import (
	"errors"
	"net/http"

	"tailings-dam/internal/model"
	"tailings-dam/internal/service"
)

// HandleListDams 列出所有尾矿坝
func (h *Handler) HandleListDams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	if status != "" {
		dams, err := h.damService.ListDamsByStatus(ctx, model.DamStatus(status))
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, dams)
		return
	}

	dams, err := h.damService.ListDams(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dams)
}

// HandleCreateDam 创建尾矿坝
func (h *Handler) HandleCreateDam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.DamInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dam, err := h.damService.CreateDam(ctx, &input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, dam)
}

// HandleGetDam 获取尾矿坝详情
func (h *Handler) HandleGetDam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	dam, err := h.damService.GetDam(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrDamNotFound) {
			respondError(w, http.StatusNotFound, "dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dam)
}

// HandleUpdateDam 更新尾矿坝
func (h *Handler) HandleUpdateDam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input model.DamInput
	if err := decodeJSON(r, &input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dam, err := h.damService.UpdateDam(ctx, id, &input)
	if err != nil {
		if errors.Is(err, service.ErrDamNotFound) {
			respondError(w, http.StatusNotFound, "dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, dam)
}

// HandleDeleteDam 删除尾矿坝
func (h *Handler) HandleDeleteDam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.damService.DeleteDam(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrDamNotFound) {
			respondError(w, http.StatusNotFound, "dam not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "dam deleted"})
}
