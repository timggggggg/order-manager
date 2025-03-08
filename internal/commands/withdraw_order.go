package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type WithdrawStorage interface {
	WithdrawOrder(ctx context.Context, id int64) (*models.Order, error)
}

type WithdrawOrder struct {
	strg WithdrawStorage
}

func NewWithdrawOrder(strg WithdrawStorage) *WithdrawOrder {
	return &WithdrawOrder{strg}
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

	w.WriteHeader(http.StatusOK)
}
