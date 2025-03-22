package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type IssueStorage interface {
	IssueOrders(ctx context.Context, ids []int64) (models.OrdersSliceStorage, error)
}

type IssueOrder struct {
	strg        IssueStorage
	logPipeline ILogPipeline
}

func NewIssueOrder(strg IssueStorage, logPipeline ILogPipeline) *IssueOrder {
	return &IssueOrder{strg, logPipeline}
}

func (cmd *IssueOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var ids []int64
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	orders, err := cmd.strg.IssueOrders(ctx, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		return
	}

	for _, order := range orders {
		cmd.logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusIssued)
	}
}
