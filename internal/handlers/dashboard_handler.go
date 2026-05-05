package handlers

import (
	"log/slog"
	"net/http"

	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type DashboardHandler struct {
	dashboardStore store.DashboardStore
	logger         *slog.Logger
}

func NewDashboardHandler(dashboardStore store.DashboardStore, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{dashboardStore: dashboardStore, logger: logger}
}

func (dh *DashboardHandler) HandleGetRetailerDashboard(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, dh.logger, "get retailer dashboard", err)
		return
	}

	dashboard, err := dh.dashboardStore.GetRetailerDashboard(id)
	if err != nil {
		utils.ServerError(w, dh.logger, "get retailer dashboard", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "dashboard fetched successfully", "dashboard": dashboard})
}
