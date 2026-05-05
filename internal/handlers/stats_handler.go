package handlers

import (
	"log/slog"
	"net/http"

	"github.com/levionstudio/fintech/internal/store"
	"github.com/levionstudio/fintech/internal/utils"
)

type StatsHandler struct {
	statsStore store.StatsStore
	logger     *slog.Logger
}

func NewStatsHandler(statsStore store.StatsStore, logger *slog.Logger) *StatsHandler {
	return &StatsHandler{statsStore: statsStore, logger: logger}
}

func (sh *StatsHandler) HandleGetRetailerStats(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, sh.logger, "get retailer stats", err)
		return
	}

	stats, err := sh.statsStore.GetRetailerStats(id)
	if err != nil {
		utils.ServerError(w, sh.logger, "get retailer stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "stats fetched successfully", "stats": stats})
}

func (sh *StatsHandler) HandleGetDistributorStats(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, sh.logger, "get distributor stats", err)
		return
	}

	stats, err := sh.statsStore.GetDistributorStats(id)
	if err != nil {
		utils.ServerError(w, sh.logger, "get distributor stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "stats fetched successfully", "stats": stats})
}

func (sh *StatsHandler) HandleGetMasterDistributorStats(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, sh.logger, "get master distributor stats", err)
		return
	}

	stats, err := sh.statsStore.GetMasterDistributorStats(id)
	if err != nil {
		utils.ServerError(w, sh.logger, "get master distributor stats", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "stats fetched successfully", "stats": stats})
}
