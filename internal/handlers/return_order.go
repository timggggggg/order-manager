package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
	logpipeline "gitlab.ozon.dev/timofey15g/homework/log_pipeline"
)

type ReturnStorage interface {
	ReturnOrder(ctx context.Context, orderID int64, userID int64) (order *models.Order, err error)
}

type ReturnOrder struct {
	strg ReturnStorage
}

func NewReturnOrder(strg ReturnStorage) *ReturnOrder {
	return &ReturnOrder{strg}
}

func (cmd *ReturnOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logPipeline := logpipeline.GetLogPipelineInstance()

	orderIDstr := r.URL.Query().Get("order_id")
	if orderIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userIDstr := r.URL.Query().Get("user_id")
	if userIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := cmd.strg.ReturnOrder(ctx, orderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}

	logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusIssued, models.StatusReturned)
}
