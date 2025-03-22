package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type WithdrawStorage interface {
	WithdrawOrder(ctx context.Context, id int64) (*models.Order, error)
}

type WithdrawOrder struct {
	strg        WithdrawStorage
	logPipeline ILogPipeline
}

func NewWithdrawOrder(strg WithdrawStorage, logPipeline ILogPipeline) *WithdrawOrder {
	return &WithdrawOrder{strg, logPipeline}
}

func (cmd *WithdrawOrder) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderIDstr := r.URL.Query().Get("order_id")
	if orderIDstr == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := cmd.strg.WithdrawOrder(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}

	cmd.logPipeline.LogStatusChange(time.Now(), order.ID, models.StatusAccepted, models.StatusWithdrawed)
}
